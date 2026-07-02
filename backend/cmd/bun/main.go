package main

import (
	"context"
	"os"

	"backend/migrations"
	"backend/internal/database"

	"github.com/uptrace/bun/migrate"
	migratorCli "github.com/iambpn/bun-migrator-cli"
)

func main() {
	db := database.InitDatabase()
	if err := db.PingContext(context.Background()); err != nil {
		os.Exit(1)
	}

	migrator := migrate.NewMigrator(db, migrations.Migrations)
	migratorCli.InitCli(migrator, os.Args)
}