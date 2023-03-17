package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	App
}

type App struct {
	GoEnv           string
	TgToken         string
	KucoinApiKey    string
	KucoinApiSecret string
	DbUrl           string
}

func GetEnv(key, fallback string) string {
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("[ERROR] loading .env file", err)
		}
	}

	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func DefaultConfig() *Config {
	var c Config

	//KG: settings App
	c.GoEnv = GetEnv("GO_ENV", "development")
	c.TgToken = GetEnv("TG_TOKEN", "")
	c.KucoinApiKey = GetEnv("KUCOIN_API_KEY", "")
	c.KucoinApiSecret = GetEnv("KUCOIN_API_SECRET", "")

	c.DbUrl = GetEnv("DATABASE_URL", "")

	return &c
}
