package migrations

import (
	"context"
	"fmt"

	"backend/internal/models/rounds"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		_, err := db.NewCreateTable().
			IfNotExists().
			Model((*rounds.Round)(nil)).
			ForeignKey(`("game_id") REFERENCES "game" ("id")`).
			ForeignKey(`("loc_id") REFERENCES "locations" ("location_id")`).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDropTable().
			IfExists().
			Model((*rounds.Round)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
