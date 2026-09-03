package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"backend/internal/models/games"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		_, err := db.NewCreateTable().
			IfNotExists().
			Model((*games.Lobby)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDropTable().
			IfExists().
			Model((*games.Lobby)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
