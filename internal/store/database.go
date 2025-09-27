package store

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"

	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
)

func Open() (*sql.DB, *redis.Client, error) {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf("db:open: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, nil, fmt.Errorf("redis:ping: %w", err)
	}
	log.Println("db:open: connected to database")
	log.Println("redis:open: connected to redis")
	return db, client, nil
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
