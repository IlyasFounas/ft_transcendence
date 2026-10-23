package migrations

import (
	"context"
	"fmt"

	"backend/internal/models/guesses"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		_, err := db.NewCreateTable().
			IfNotExists().
			Model((*guesses.Guess)(nil)).
			WithForeignKeys().
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDropTable().
			IfExists().
			Model((*guesses.Guess)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
