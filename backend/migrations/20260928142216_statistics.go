package migrations

import (
	"context"
	"fmt"

	"backend/internal/models/users"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		_, err := db.NewCreateTable().
			IfNotExists().
			Model((*users.Statistics)(nil)).
			ForeignKey(`("user_id") REFERENCES "users" ("id") ON DELETE CASCADE`).
			Exec(ctx)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDropTable().
			IfExists().
			Model((*users.Statistics)(nil)).
			Exec(ctx)
		return err
	})
}
