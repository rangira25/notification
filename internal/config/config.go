package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Env               string
	RedisAddr         string
	RedisPassword     string
	DatabaseURL       string
	APIAddr           string
	WorkerConcurrency int
	AsynqmonPort      string
	SMTPHost          string
	SMTPPort          string
	SMTPUser          string
	SMTPPass          string
}

func Load() *Config {
	c := &Config{
		Env:           getenv("ENV", "development"),
		RedisAddr:     getenv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getenv("REDIS_PASSWORD", ""),
		DatabaseURL:   getenv("DATABASE_URL", ""),
		APIAddr:       getenv("API_ADDR", ":8080"),
		AsynqmonPort:  getenv("ASYNQMON_PORT", "8081"),
		SMTPHost:      getenv("SMTP_HOST", ""),
		SMTPPort:      getenv("SMTP_PORT", "587"),
		SMTPUser:      getenv("SMTP_USER", ""),
		SMTPPass:      getenv("SMTP_PASS", ""),
	}

	wc := getenv("WORKER_CONCURRENCY", "10")
	n, err := strconv.Atoi(wc)
	if err != nil {
		n = 10
	}
	c.WorkerConcurrency = n

	if c.DatabaseURL == "" {
		log.Fatal("DATABASE_URL not set")
	}
	return c
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
