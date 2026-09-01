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
			Model((*games.Game)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		
		// Add CHECK constraints
		_, err = db.Exec(`ALTER TABLE "game" ADD CONSTRAINT chk_round_total CHECK (round_total <> 0)`)
		if err != nil && !isConstraintExistsError(err) {
			return err
		}
		
		_, err = db.Exec(`ALTER TABLE "game" ADD CONSTRAINT chk_player_nb CHECK (player_nb <> 0)`)
		if err != nil && !isConstraintExistsError(err) {
			return err
		}
		
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDropTable().
			IfExists().
			Model((*games.Game)(nil)).
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func isConstraintExistsError(err error) bool {
	// Ignore "constraint already exists" errors
	return err != nil && (err.Error() == "pq: constraint \"chk_round_total\" for relation \"game\" already exists" ||
		err.Error() == "pq: constraint \"chk_player_nb\" for relation \"game\" already exists")
}
