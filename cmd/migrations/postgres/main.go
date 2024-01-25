package main

import (
	"database/sql"
	"effective-mobile-test-task/internal/lib/config"
	"effective-mobile-test-task/internal/repositories/postgres"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {
	// Config load
	envPath := os.Getenv("ENV_PATH")

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("failed to load env file. Err: %s", err.Error())
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to create config. Err: %s", err.Error())
	}

	// PG init
	pg, err := postgres.New(cfg.PG)
	if err != nil {
		log.Fatalf("failed to create database instance. Err: %s", err.Error())
	}
	defer pg.Close()

	// Getting steps
	if len(os.Args) == 1 {
		log.Fatalf("operation not specified")
	}

	// Migrating
	switch os.Args[1] {
	case "top":
		err = migrateTop(pg, cfg.PG.MigrationsPath)
	case "drop":
		err = dropMigrations(pg, cfg.PG.MigrationsPath)
	default:
		err = migrateNSteps(pg, cfg.PG.MigrationsPath, os.Args[1])
	}

	if err != nil {
		log.Fatalf("error occuired: %s", err.Error())
	}
}

func migrateTop(pg *sql.DB, migrationsPath string) error {
	return postgres.MigrateTop(pg, migrationsPath)
}

func dropMigrations(pg *sql.DB, migrationsPath string) error {
	return postgres.DropMigrations(pg, migrationsPath)
}

func migrateNSteps(pg *sql.DB, migrationsPath string, n string) error {
	steps, err := strconv.Atoi(n)
	if err != nil {
		return fmt.Errorf("wrong type of argument: %w", err)
	}
	if steps == 0 {
		return fmt.Errorf("wrong number of steps. n > 0 to migrate up and n < 0 to migrate down")
	}

	return postgres.MigrateNSteps(pg, migrationsPath, steps)
}
