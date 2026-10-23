package migrations

import (
	"backend/internal/models/friends"
	"context"
	"fmt"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		_, err := db.NewCreateTable().
			IfNotExists().
			Model((*friends.Friend)(nil)).
			ForeignKey(`("requester_id") REFERENCES "users" ("id") ON DELETE CASCADE`).
			ForeignKey(`("receiver_id") REFERENCES "users" ("id") ON DELETE CASCADE`).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDropTable().
			IfExists().
			Model((*friends.Friend)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
