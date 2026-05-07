package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"rdev-go-api/internal/config"
	"rdev-go-api/internal/data"

	"github.com/uptrace/bun/migrate"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	db := data.ConnectDB(cfg.Database.URL, cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns)
	migrator := migrate.NewMigrator(db, data.Migrations)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: go run cmd/migrate/main.go [init|up|down|status]")
		return
	}

	cmd := args[1]
	switch strings.ToLower(cmd) {
	case "init":
		err = migrator.Init(ctx)
	case "up":
		group, err := migrator.Migrate(ctx)
		if err == nil {
			fmt.Printf("Migrated to %s\n", group)
		}
	case "down":
		group, err := migrator.Rollback(ctx)
		if err == nil {
			fmt.Printf("Rolled back %s\n", group)
		}
	case "status":
		ms, err := migrator.MigrationsWithStatus(ctx)
		if err == nil {
			fmt.Printf("Migration Status:\n%s\n", ms)
		}
	}

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
