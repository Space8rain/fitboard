package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Env string
    BotToken   string
    DBConnStr  string
    WebAppURL  string
}

func Load() *Config {
    err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

    cfg := &Config{
        Env:       os.Getenv("APP_ENV"),
        BotToken:  os.Getenv("TOKEN_BOT"),
        DBConnStr: os.Getenv("DB_CONN_STR"),
        WebAppURL: os.Getenv("WEB_APP_URL"),
    }

    if cfg.Env == "" {
        log.Fatal("APP_ENV is required")
    }

    if cfg.BotToken == "" {
        log.Fatal("TOKEN_BOT is required")
    }

    if cfg.DBConnStr == "" {
        log.Fatal("DB_CONN_STR is required")
    }

    if cfg.WebAppURL == "" {
        log.Fatal("WEB_APP_URL is required")
    }

    return cfg
}
