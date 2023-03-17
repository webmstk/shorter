package config

import (
	"flag"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
)

type AppConfig struct {
	ServerAddress    string `env:"SERVER_ADDRESS"`
	BaseURL          string `env:"BASE_URL"`
	Schema           string
	FileStoragePath  string `env:"FILE_STORAGE_PATH"`
	CookieSalt       string `env:"COOKIE_SALT" envDefault:"ABRVAL6"`
	CookieTTLSeconds int    `env:"COOKIE_TTL_SECONDS" envDefault:"604800"`
}

var Config AppConfig
var ServerBaseURL string

func InitConfig() {
	Config.Schema = "http"

	parseCliArguments()
	parseEnvVariables()

	ServerBaseURL = fmt.Sprintf("%s://%s", Config.Schema, Config.ServerAddress)
}

func parseCliArguments() {
	serverAddress := flag.String("a", "localhost:8080", "server address")
	baseURL := flag.String("b", "http://localhost:8080", "base server address")
	fileStoragePath := flag.String("f", "storage/storage.json", "file storage path")

	flag.Parse()

	Config.ServerAddress = *serverAddress
	Config.BaseURL = *baseURL
	Config.FileStoragePath = *fileStoragePath
}

func parseEnvVariables() {
	err := env.Parse(&Config)
	if err != nil {
		log.Fatal(err)
	}
}
