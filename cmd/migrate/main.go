package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/XaiPhyr/rdev-go-api/internal/config"
	"github.com/XaiPhyr/rdev-go-api/internal/db"

	"github.com/uptrace/bun/migrate"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	args := os.Args
	if len(args) < 2 {
		log.Println("Usage: go run cmd/migrate/main.go [init|up|down|status]")
		return
	}

	cmd := args[1]

	if err := CheckArgs(cmd, cfg.Database); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

func CheckArgs(cmd string, dbConfig config.DBConfig) error {
	var err error
	ctx := context.Background()

	connDB := config.ConnectDB(dbConfig)
	migrator := migrate.NewMigrator(connDB, db.Migrations)

	switch strings.ToLower(cmd) {
	case "init":
		err = migrator.Init(ctx)
	case "up":
		group, err := migrator.Migrate(ctx)
		if err == nil {
			log.Printf("Migrated to %s\n", group)
		}
	case "down":
		group, err := migrator.Rollback(ctx)
		if err == nil {
			log.Printf("Rolled back %s\n", group)
		}
	case "status":
		ms, err := migrator.MigrationsWithStatus(ctx)
		if err == nil {
			log.Printf("Migration Status:\n%s\n", ms)
		}
	}

	return err
}
