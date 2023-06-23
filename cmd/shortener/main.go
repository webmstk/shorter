package main

import (
	"log"

	"github.com/webmstk/shorter/internal/app"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/db"
)

func main() {
	config.InitConfig()
	if config.Config.DatabaseDSN != "" {
		err := db.Init()
		if err != nil {
			log.Fatal(err)
		}
	}
	app.Run()
}
