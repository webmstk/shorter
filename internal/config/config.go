package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
)

type AppConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	Schema        string
}

var Config AppConfig
var ServerBaseURL string

func init() {
	Config.Schema = "http"

	err := env.Parse(&Config)
	if err != nil {
		log.Fatal(err)
	}

	ServerBaseURL = fmt.Sprintf("%s://%s", Config.Schema, Config.ServerAddress)
}
