package tests

import (
	"context"
	"flag"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/db"
)

func Setup() {
	setupTestConfig(&config.Config)
}

func setupTestConfig(config *config.AppConfig) {
	config.ServerAddress = "localhost:8080"
	config.BaseURL = "http://localhost:8080"
	config.FileStoragePath = ""
	config.CookieTTLSeconds = 120
	config.CookieSalt = "secret"

	withDB := flag.Bool("with-db", false, "database DSN")
	flag.Parse()

	if *withDB {
		config.DatabaseDSN = "postgres://postgres:@localhost:5432/shorter_test?sslmode=disable"
		db.Init()
		// ClearTables()
	}
}

func ClearTables() {
	conn, err := pgx.Connect(context.Background(), config.Config.DatabaseDSN)

	if err != nil {
		log.Fatal(err)
	}

	sql := ""

	tables := []string{"user_links", "users", "links"}
	for _, table := range tables {
		sql = sql + "TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE;\n"
	}

	_, err = conn.Exec(context.Background(), sql)
	if err != nil {
		log.Fatal(err)
	}
}
