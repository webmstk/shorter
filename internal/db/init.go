package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/webmstk/shorter/internal/config"
)

func Init() {
	databaseDSN := config.Config.DatabaseDSN
	conn, err := pgx.Connect(context.Background(), databaseDSN)
	if err != nil {
		log.Fatal(err)
	}

	createUsers := `CREATE TABLE IF NOT EXISTS users (
	  uuid VARCHAR ( 36 ) UNIQUE NOT NULL,
	  created_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`

	createLinks := `CREATE TABLE IF NOT EXISTS links (
          id SERIAL PRIMARY KEY,
          short_url VARCHAR UNIQUE NOT NULL,
	  long_url VARCHAR UNIQUE NOT NULL,
	  created_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`

	createUserLinks := `CREATE TABLE IF NOT EXISTS user_links (
          user_id VARCHAR ( 36 ) NOT NULL,
	  link_id INT NOT NULL,
	  FOREIGN KEY (user_id) REFERENCES users (uuid),
	  FOREIGN KEY (link_id) REFERENCES links (id)

	)`

	createUserLinksIndex := `CREATE UNIQUE INDEX IF NOT EXISTS userid_linkid ON user_links(user_id, link_id)`

	runQuery(conn, createUsers)
	runQuery(conn, createLinks)

	runQuery(conn, createUserLinks)
	runQuery(conn, createUserLinksIndex)
}

func runQuery(conn *pgx.Conn, sql string) {
	_, err := conn.Exec(context.Background(), sql)
	if err != nil {
		log.Fatal(err)
	}
}
