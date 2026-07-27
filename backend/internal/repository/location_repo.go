package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/uptrace/bun"

	"backend/internal/model"
)

var ErrLocationNotFound = errors.New("location not found")

// LocationRepository is a repo for managing locations in the database and synchronizing with a JSON file.
type LocationRepository struct {
	db       *bun.DB
	filePath string
}

func NewLocationRepository(db *bun.DB, filePath string) *LocationRepository {
	return &LocationRepository{db: db, filePath: filePath}
}

// syncFile re-reads the entire table and rewrites the JSON file accordingly.
// Called after each successful mutation (create/patch/delete).
func (r *LocationRepository) syncFile(ctx context.Context) error {
	locs, err := r.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("reading locations for file sync: %w", err)
	}

	data, err := json.MarshalIndent(locs, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling locations: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("writing locations file: %w", err)
	}
	return nil
}

// GetAll — all locations used internally for synchronization and exposed if needed
func (r *LocationRepository) GetAll(ctx context.Context) ([]*model.Location, error) {
	var locs []*model.Location
	err := r.db.NewSelect().Model(&locs).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return locs, nil
}

// CREATE - add multiple locations to the database, ignores duplicates based on Mapillary ID. Returns the number of new locations added.
func (r *LocationRepository) BulkCreate(ctx context.Context, locations []*model.Location) (int, error) {
	if len(locations) == 0 {
		return 0, nil
	}

	res, err := r.db.NewInsert().
		Model(&locations).
		On("CONFLICT (id_mapillary) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if affected > 0 {
		if err := r.syncFile(ctx); err != nil {
			return int(affected), err // means insertion in DB succeeded, but file synchronization failed
		}
	}

	return int(affected), nil
}

// READ (using internal ID)
func (r *LocationRepository) GetByID(ctx context.Context, id int64) (*model.Location, error) {
	loc := new(model.Location)
	err := r.db.NewSelect().Model(loc).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return loc, nil
}

// READ (using Mapillary ID)
func (r *LocationRepository) GetByMapillaryID(ctx context.Context, mid int64) (*model.Location, error) {
	loc := new(model.Location)
	err := r.db.NewSelect().Model(loc).Where("id_mapillary = ?", mid).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return loc, nil
}

// READ - random location
func (r *LocationRepository) GetRandom(ctx context.Context) (*model.Location, error) {
	loc := new(model.Location)
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
func (r *LocationRepository) GetRandomByCountryCode(ctx context.Context, cc string) (*model.Location, error) {
	loc := new(model.Location)
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
func (r *LocationRepository) PatchByID(ctx context.Context, id int64, patch *model.Location) error {
	return r.patch(ctx, "id = ?", id, patch)
}

// PATCH (using Mapillary ID)
func (r *LocationRepository) PatchByMapillaryID(ctx context.Context, mid int64, patch *model.Location) error {
	return r.patch(ctx, "id_mapillary = ?", mid, patch)
}

// patch — private helper function to update a location based on a where clause and patch data. Only updates fields that are not null.
func (r *LocationRepository) patch(ctx context.Context, whereClause string, whereArg any, patch *model.Location) error {
	query := r.db.NewUpdate().
		Model((*model.Location)(nil)).
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

	return r.syncFile(ctx)
}

// DELETE (using internal ID)
func (r *LocationRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.delete(ctx, "id = ?", id)
}

// DELETE (using Mapillary ID)
func (r *LocationRepository) DeleteByMapillaryID(ctx context.Context, mid int64) error {
	return r.delete(ctx, "id_mapillary = ?", mid)
}

// delete — private helper function for common deletion logic
func (r *LocationRepository) delete(ctx context.Context, whereClause string, whereArg any) error {
	res, err := r.db.NewDelete().
		Model((*model.Location)(nil)).
		Where(whereClause, whereArg).
		Exec(ctx)
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

	return r.syncFile(ctx)
}