package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr        string
	DatabaseURL string
}

func FromEnv() Config {
	loadEnvFiles()

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":18080"
	}

	return Config{
		Addr:        addr,
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}

func loadEnvFiles() {
	for _, path := range []string{".env.local", ".env", "../.env.local", "../.env"} {
		_ = godotenv.Load(path)
	}
}
