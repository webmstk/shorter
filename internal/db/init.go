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

	log.Println("--- CREATING ---")
	log.Println("users table creating")
	runQuery(conn, createUsers)
	log.Println("users table created")
	log.Println("links table creating")
	runQuery(conn, createLinks)
	log.Println("links table created")

	log.Println("user_links table creating")
	runQuery(conn, createUserLinks)
	log.Println("user_links table created")
	log.Println("user_links index creating")
	runQuery(conn, createUserLinksIndex)
	log.Println("user_links index created")

	log.Println("--- DONE ---")
}

func runQuery(conn *pgx.Conn, sql string) {
	_, err := conn.Exec(context.Background(), sql)
	if err != nil {
		log.Fatal(err)
	}
}
