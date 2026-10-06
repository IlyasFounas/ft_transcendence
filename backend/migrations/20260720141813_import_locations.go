package migrations

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/uptrace/bun"

	"backend/internal/models/locations"
)

// next line is a go directive (go compiler embeds the contents of the file into the binary)
//go:embed data/locations.json
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

		locationList := make([]*locations.Location, 0, len(imports))
		for _, item := range imports {
			if len(item.CountryCode) != 2 {
				continue
			}
			locationList = append(locationList, &locations.Location{
				MID:  &item.MapillaryID,
				Lat:  &item.Latitude,
				Long: &item.Longitude,
				CC:   &item.CountryCode,
			})
		}

		if len(locationList) == 0 {
			return nil
		}

		_, err = db.NewInsert().
			Model(&locationList).
			On("CONFLICT (mapillary_id) DO NOTHING").
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
		location_ids := make([]int64, 0, len(imports))
		for _, item := range imports {
			location_ids = append(location_ids, item.MapillaryID)
		}
		_, err = db.NewDelete().
			Model((*locations.Location)(nil)).
			Where("mapillary_id IN ?", bun.Tuple(location_ids)).
			Exec(ctx)
		return err
	})
}
