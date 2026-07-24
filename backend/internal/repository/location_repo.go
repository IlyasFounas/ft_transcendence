// TODO: SELECT * FROM table ORDER BY RANDOM() LIMIT 1

package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
	"backend/internal/models"
)

var ErrLocationNotFound = errors.New("location not found")

type LocationRepository struct {
	db *bun.DB
}

func NewLocationRepository(db *bun.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

// CREATE - add multiple locations to the database.
// If a location with the same Mapillary ID already exists, it is ignored.
func (r *LocationRepository) BulkCreate(ctx context.Context, locations []*models.Location) (int, error) {
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
	return int(affected), nil
}

// READ (using ID)
func (r *LocationRepository) GetByID(ctx context.Context, id int64) (*models.Location, error) {
	loc := new(models.Location)
	err := r.db.NewSelect().
		Model(loc).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return loc, nil
}

// READ (using Mapillary ID)
func (r *LocationRepository) GetByMapillaryID(ctx context.Context, mid int64) (*models.Location, error) {
	loc := new(models.Location)
	err := r.db.NewSelect().
		Model(loc).
		Where("id_mapillary = ?", mid).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return loc, nil
}

// READ - random location
func (r *LocationRepository) GetRandom(ctx context.Context) (*models.Location, error) {
	loc := new(models.Location)
	err := r.db.NewSelect().
		Model(loc).
		OrderExpr("RANDOM()").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return loc, nil
}

// READ - random location filtered by country code
func (r *LocationRepository) GetRandomByCountryCode(ctx context.Context, cc string) (*models.Location, error) {
	loc := new(models.Location)
	err := r.db.NewSelect().
		Model(loc).
		Where("country_code = ?", cc).
		OrderExpr("RANDOM()").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return loc, nil
}

// PATCH - update a location by its internal ID
func (r *LocationRepository) PatchByID(ctx context.Context, id int64, patch models.LocationPatch) error {
	return r.patch(ctx, "id = ?", id, patch)
}

// PATCH - update a location by its Mapillary ID
func (r *LocationRepository) PatchByMapillaryID(ctx context.Context, mid int64, patch models.LocationPatch) error {
	return r.patch(ctx, "id_mapillary = ?", mid, patch)
}

// PATCH - private helper function to update a location based on a where clause and patch data
func (r *LocationRepository) patch(ctx context.Context, whereClause string, whereArg any, patch models.LocationPatch) error {
	query := r.db.NewUpdate().
		Model((*models.Location)(nil)).
		Where(whereClause, whereArg)

	hasChange := false
	if patch.Latitude != nil {
		query.Set("latitude = ?", *patch.Latitude)
		hasChange = true
	}
	if patch.Longitude != nil {
		query.Set("longitude = ?", *patch.Longitude)
		hasChange = true
	}
	if patch.CountryCode != nil {
		query.Set("country_code = ?", *patch.CountryCode)
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
	return nil
}

// DELETE - delete a location by its internal ID
func (r *LocationRepository) DeleteByID(ctx context.Context, id int64) error {
	res, err := r.db.NewDelete().
		Model((*models.Location)(nil)).
		Where("id = ?", id).
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
	return nil
}

// DELETE - delete a location by its Mapillary ID
func (r *LocationRepository) DeleteByMapillaryID(ctx context.Context, mid int64) error {
	res, err := r.db.NewDelete().
		Model((*models.Location)(nil)).
		Where("id_mapillary = ?", mid).
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
	return nil
}