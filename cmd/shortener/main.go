package main

import (
	"github.com/webmstk/shorter/internal/app"
	"github.com/webmstk/shorter/internal/config"
)

func main() {
	config.InitConfig()
	app.Run()
}
