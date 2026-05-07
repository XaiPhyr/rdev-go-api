package data

import (
	"database/sql"
	"log"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func ConnectDB(dsn string, maxOpenConns, maxIdleConns int) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxIdleConns)
	sqldb.SetConnMaxLifetime(1 * time.Hour)
	sqldb.SetConnMaxIdleTime(30 * time.Minute)

	db := bun.NewDB(sqldb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(true),
		bundebug.WithVerbose(true),
	))

	if err := db.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
		return nil
	}

	return db
}
