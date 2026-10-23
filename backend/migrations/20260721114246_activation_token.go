package migrations

import (
	"backend/internal/models/users/tokens"
	"context"
	"fmt"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		_, err := db.NewCreateTable().
			IfNotExists().
			Model((*tokens.ActivationToken)(nil)).
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
			Model((*tokens.ActivationToken)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
