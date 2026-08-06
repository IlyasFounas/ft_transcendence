package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"backend/internal/models/locations"
	"github.com/uptrace/bun"
)

var ErrLocationNotFound = errors.New("location not found")
var ErrLocationNotSet = errors.New("location path env variable not set")

// LocationRepository is a repo for managing locations in the database and synchronizing with a JSON file.
type LocationRepository struct {
	db         *bun.DB
	logPath    string
	logChannel chan LogEvent
	wg         sync.WaitGroup // Enables to wait that the worker finished to write everything
}

func NewLocationRepository(db *bun.DB) (*LocationRepository, error) {
	logPath := os.Getenv("LOCATION_FILE_PATH")
	if logPath == "" {
		return nil, ErrLocationNotSet
	}

	repo := &LocationRepository{
		db:         db,
		logPath:    logPath,
		logChannel: make(chan LogEvent, 200),
	}

	repo.wg.Add(1) // Signals that a goroutine will be started
	go repo.startLogWorker()

	return repo, nil
}

type LogEvent struct {
	Timestamp time.Time           `json:"timestamp"`
	Action    string              `json:"action"` // CREATE, PATCH, DELETE
	Location  *locations.Location `json:"location,omitempty"`
}

// Method to call when shutting down the server
func (r *LocationRepository) Close() {
	close(r.logChannel) // 1. Disable new sends and signal the end of the stream
	r.wg.Wait()         // 2. Block until all remaining events are written
}

func (r *LocationRepository) startLogWorker() {
	defer r.wg.Done() // Indicate that the worker has finished when it exits the loop

	// 'range' stops automatically as soon as the channel is closed (in Close())
	for event := range r.logChannel {
		if err := r.appendJSONL(event); err != nil {
			fmt.Printf("Error writing to JSONL: %v\n", err)
		}
	}
}

// appendJSONL adds an event in JSONL format and creates the parent directory if needed
func (r *LocationRepository) appendJSONL(event LogEvent) error {
	// Ensures the parent directory exists
	if err := os.MkdirAll(filepath.Dir(r.logPath), 0755); err != nil {
		return fmt.Errorf("creating log directory: %w", err)
	}

	f, err := os.OpenFile(r.logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening log file: %w", err)
	}
	defer f.Close()

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshaling log event: %w", err)
	}

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("writing to log file: %w", err)
	}
	return nil
}

// CREATE - add multiple locations to the database, ignores duplicates based on Mapillary ID. Returns the number of new locations added.
func (r *LocationRepository) BulkCreate(ctx context.Context, locations []*locations.Location) (int, error) {
	if len(locations) == 0 {
		return 0, nil
	}

	res, err := r.db.NewInsert().
		Model(&locations).
		On("CONFLICT (mapillary_id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if affected > 0 {
		// loop of logging for each created location
		for _, loc := range locations {
			r.logChannel <- LogEvent{
				Timestamp: time.Now().UTC(),
				Action:    "CREATE",
				Location:  loc,
			}
		}
	}
	return int(affected), nil
}

// READ (using internal ID)
func (r *LocationRepository) GetByID(ctx context.Context, location_id int64) (*locations.Location, error) {
	loc := new(locations.Location)
	err := r.db.NewSelect().Model(loc).Where("location_id = ?", location_id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return loc, nil
}

// READ (using Mapillary ID)
func (r *LocationRepository) GetByMapillaryID(ctx context.Context, mid int64) (*locations.Location, error) {
	loc := new(locations.Location)
	err := r.db.NewSelect().Model(loc).Where("mapillary_id = ?", mid).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return loc, nil
}

// READ - random location
func (r *LocationRepository) GetRandom(ctx context.Context) (*locations.Location, error) {
	loc := new(locations.Location)
	err := r.db.NewSelect().Model(loc).OrderExpr("RANDOM()").Limit(1).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return loc, nil
}

// READ — random location filtered by country code
func (r *LocationRepository) GetRandomByCountryCode(ctx context.Context, cc string) (*locations.Location, error) {
	loc := new(locations.Location)
	err := r.db.NewSelect().Model(loc).Where("country_code = ?", cc).OrderExpr("RANDOM()").Limit(1).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return loc, nil
}

// PATCH (using internal ID)
func (r *LocationRepository) PatchByID(ctx context.Context, location_id int64, patch *locations.Location) error {
	return r.patch(ctx, "location_id = ?", location_id, patch)
}

// PATCH (using Mapillary ID)
func (r *LocationRepository) PatchByMapillaryID(ctx context.Context, mid int64, patch *locations.Location) error {
	return r.patch(ctx, "mapillary_id = ?", mid, patch)
}

// patch — private helper function to update a location based on a where clause and patch data. Only updates fields that are not null.
func (r *LocationRepository) patch(ctx context.Context, whereClause string, whereArg any, patch *locations.Location) error {
	query := r.db.NewUpdate().
		Model((*locations.Location)(nil)).
		Where(whereClause, whereArg)

	hasChange := false
	if patch.Lat != nil {
		query.Set("latitude = ?", *patch.Lat)
		hasChange = true
	}
	if patch.Long != nil {
		query.Set("longitude = ?", *patch.Long)
		hasChange = true
	}
	if patch.CC != nil {
		query.Set("country_code = ?", *patch.CC)
		hasChange = true
	}

	if !hasChange {
		return nil
	}

	res, err := query.Exec(ctx)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrLocationNotFound
	}

	r.logChannel <- LogEvent{
		Timestamp: time.Now().UTC(),
		Action:    "PATCH",
		Location:  patch,
	}

	return nil
}

// DELETE (using internal ID)
func (r *LocationRepository) DeleteByID(ctx context.Context, location_id int64) error {
	return r.delete(ctx, "location_id = ?", location_id)
}

// DELETE (using Mapillary ID)
func (r *LocationRepository) DeleteByMapillaryID(ctx context.Context, mid int64) error {
	return r.delete(ctx, "mapillary_id = ?", mid)
}

// delete — private helper function for common deletion logic
func (r *LocationRepository) delete(ctx context.Context, whereClause string, whereArg any) error {
	//Save the location before deletion for logging
	loc := new(locations.Location)
	err := r.db.NewSelect().
		Model(loc).
		Where(whereClause, whereArg).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrLocationNotFound
		}
		return err
	}
	// Perform the deletion
	res, err := r.db.NewDelete().
		Model((*locations.Location)(nil)).
		Where(whereClause, whereArg).
		Exec(ctx)
	if err != nil {
		return err
	}

	// Check if any row was deleted
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrLocationNotFound
	}

	// Log the deletion
	r.logChannel <- LogEvent{
		Timestamp: time.Now().UTC(),
		Action:    "DELETE",
		Location:  loc,
	}
	return nil
}
