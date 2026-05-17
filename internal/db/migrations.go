package db

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

//go:embed migrations/*.sql
var migrations embed.FS

var Migrations = migrate.NewMigrations()

func init() {
	if err := Migrations.Discover(migrations); err != nil {
		panic(err)
	}
}
