package migrations

import (
	"context"
	"encoding/json"
	"fmt"
	"embed"
	"github.com/uptrace/bun"

	"backend/internal/model"
)

var locationsFile embed.FS

type locationImport struct {
	MapillaryID int64   `json:"mapillary_id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	CountryCode string  `json:"country_code"`
}

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		data, err := locationsFile.ReadFile("data/locations.json")
		if err != nil {
			return fmt.Errorf("reading seed file: %w", err)
		}

		var imports []locationImport
		if err := json.Unmarshal(data, &imports); err != nil {
			return fmt.Errorf("parsing seed JSON: %w", err)
		}

		locations := make([]*model.Location, 0, len(imports))
		for _, item := range imports {
			if len(item.CountryCode) != 2 {
				continue
			}
			locations = append(locations, &model.Location{
				ID_mplry: item.MapillaryID,
				Lat:      item.Latitude,
				Long:     item.Longitude,
				CC:       item.CountryCode,
			})
		}

		if len(locations) == 0 {
			return nil
		}

		_, err = db.NewInsert().
			Model(&locations).
			On("CONFLICT (id_mapillary) DO NOTHING").
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("inserting seed locations: %w", err)
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		data, err := locationsFile.ReadFile("data/locations.json")
		if err != nil {
			return err
		}
		var imports []locationImport
		if err := json.Unmarshal(data, &imports); err != nil {
			return err
		}
		ids := make([]int64, 0, len(imports))
		for _, item := range imports {
			ids = append(ids, item.MapillaryID)
		}
		_, err = db.NewDelete().
			Model((*model.Location)(nil)).
			Where("id_mapillary IN (?)", bun.In(ids)).
			Exec(ctx)
		return err
	})
}
