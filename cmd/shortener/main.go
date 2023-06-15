package main

import (
	"github.com/webmstk/shorter/internal/app"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/db"
)

func main() {
	config.InitConfig()
	if config.Config.DatabaseDSN != "" {
		db.Init()
	}
	app.Run()
}
