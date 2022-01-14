package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	Port    string

	NATSHost string

	PGHost string
	PGUser string
	PGPass string
	PGName string
}

func New() *Config {
	_ = godotenv.Load(".env")

	c := &Config{
		AppName:  os.Getenv("APP_NAME"),
		Port:     os.Getenv("PORT"),
		NATSHost: os.Getenv("NATS_HOST"),
		PGHost:   os.Getenv("PG_HOST"),
		PGUser:   os.Getenv("PG_USER"),
		PGPass:   os.Getenv("PG_PASS"),
		PGName:   os.Getenv("PG_NAME"),
	}

	return c
}
