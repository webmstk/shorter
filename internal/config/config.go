package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
)

type AppConfig struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	Schema          string
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
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
