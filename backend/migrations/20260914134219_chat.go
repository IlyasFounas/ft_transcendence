package migrations

import (
	"backend/internal/models/chats"
	"context"
	"fmt"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		_, err := db.NewCreateTable().
			IfNotExists().
			Model((*chats.Chat)(nil)).
			ForeignKey(`("src_user_id") REFERENCES "users" ("id") ON DELETE CASCADE`).
			ForeignKey(`("dst_user_id") REFERENCES "users" ("id") ON DELETE CASCADE`).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDropTable().
			IfExists().
			Model((*chats.Chat)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
