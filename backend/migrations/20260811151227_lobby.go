package migrations

import (
	"context"
	"fmt"

	"backend/internal/models/games"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		_, err := db.NewCreateTable().
			IfNotExists().
			Model((*games.Lobby)(nil)).
			ForeignKey(`("user_id") REFERENCES "users" ("id") ON DELETE CASCADE`).
			ForeignKey(`("game_id") REFERENCES "game" ("id") ON DELETE CASCADE`).
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
