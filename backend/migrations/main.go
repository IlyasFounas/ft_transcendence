package migrations

import "github.com/uptrace/bun/migrate"

var Migrations = migrate.NewMigrations()


// Package migrations initializes and discovers database migrations
func init() {
	if err := Migrations.DiscoverCaller(); err != nil {
		panic(err)
	}
}