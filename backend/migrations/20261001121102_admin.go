package migrations

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"backend/internal/models/users"
	"backend/internal/tools/encryption"
	"github.com/uptrace/bun"
)

// Sets up admin user, and the admin account
// Its name, email and password come from the environment
// The name and the password can't contain "admin", and the password needs at least 12 characters
// If one of them is missing or refused, the migration will fail
func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		name := os.Getenv("ADMIN_USERNAME")
		email := os.Getenv("ADMIN_EMAIL")
		password := os.Getenv("ADMIN_PASSWORD")
		if name == "" || email == "" || password == "" {
			return errors.New("ADMIN_USERNAME, ADMIN_EMAIL and ADMIN_PASSWORD must be set to create the admin")
		}
		if strings.Contains(strings.ToLower(name), "admin") || strings.Contains(strings.ToLower(password), "admin") {
			return errors.New(`the name and the password of the admin can't contain "admin"`)
		}
		if len(password) < 12 {
			return errors.New("the password of the admin must have at least 12 characters")
		}
		hash, err := encryption.HashPassword(password)
		if err != nil {
			return err
		}

		// the column (does nothing if it already exists, so running the migration again is safe)
		_, err = db.NewAddColumn().
			Model((*users.User)(nil)).
			ColumnExpr("is_admin boolean NOT NULL DEFAULT false").
			IfNotExists().
			Exec(ctx)
		if err != nil {
			return err
		}

		// the account: verified (no activation email), logs in with its name as login
		admin := &users.User{
			Username:    name,
			Login:       &name,
			Email:       email,
			Password:    &hash,
			Is_verified: true,
			IsAdmin:     true,
		}
		_, err = db.NewInsert().
			Model(admin).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("could not create the admin %q (name, login or email already used?): %w", name, err)
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDelete().
			Model((*users.User)(nil)).
			Where("is_admin = ?", true).
			Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewDropColumn().
			Model((*users.User)(nil)).
			Column("is_admin").
			Exec(ctx)
		return err
	})
}
