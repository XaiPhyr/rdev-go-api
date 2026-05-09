package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/XaiPhyr/rdev-go-api/internal/config"
	"github.com/XaiPhyr/rdev-go-api/internal/data"

	"github.com/uptrace/bun/migrate"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	db := config.ConnectDB(cfg.Database)
	migrator := migrate.NewMigrator(db, data.Migrations)

	args := os.Args
	if len(args) < 2 {
		log.Println("Usage: go run cmd/migrate/main.go [init|up|down|status]")
		return
	}

	cmd := args[1]
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

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
