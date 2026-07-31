package main

import (
	"context"
	"os"

	"backend/internal/database"
	"backend/migrations"

	migratorCli "github.com/iambpn/bun-migrator-cli"
	"github.com/uptrace/bun/migrate"
)

func main() {
	db := database.InitDatabase()
	if err := db.PingContext(context.Background()); err != nil {
		os.Exit(1)
	}

	migrator := migrate.NewMigrator(db, migrations.Migrations)
	migratorCli.InitCli(migrator, os.Args)
}
