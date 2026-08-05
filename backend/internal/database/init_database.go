package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func buildDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
}

func InitDatabase() *bun.DB {
	// ctx := context.Background()

	// Specify the information for the connection
	dsn := buildDSN()
	// Connect to the PostgreSQL database
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Create a Bun db on top of it
	db := bun.NewDB(pgdb, pgdialect.New())

	// Optional: Print queries during development
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	// Ping test to check the connection with the db
	// if err := db.Ping(context.Background()); err != nil {
	//     panic(fmt.Sprintf("Impossible to connect to PostgreSQL: %v", err))
	// }
	return db
}
