package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("db:open: %w", err)
	}
	log.Println("db:open: connected to database")
	return db, nil
}
func MigrateFS(sql *sql.DB, migrationFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return migrate(sql, dir)
}

func migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("db:migrate:set-dialect: %w", err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("db:migrate:up: %w", err)
	}
	log.Println("db:migrate:up: migrations applied successfully")
	return nil
}
