package tests

import "github.com/webmstk/shorter/internal/config"

func Setup() {
	setupTestConfig(&config.Config)
}

func setupTestConfig(config *config.AppConfig) {
	config.ServerAddress = "localhost:8080"
	config.BaseURL = "http://localhost:8080"
	config.FileStoragePath = ""
	config.CookieTTLSeconds = 120
	config.CookieSalt = "secret"
	// config.DatabaseDSN = "postgres://postgres:@localhost:5432/shorter_test?sslmode=disable"
}
