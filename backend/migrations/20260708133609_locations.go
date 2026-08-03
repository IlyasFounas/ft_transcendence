package migrations

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"

	"backend/internal/model/locations"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		_, err := db.NewCreateTable().
			IfNotExists().
			Model((*locations.Location)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDropTable().
			IfExists().
			Model((*locations.Location)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
