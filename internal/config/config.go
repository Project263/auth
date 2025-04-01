package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	POSTGRES_DSN     string `env:"POSTGRES_DSN" envDefault:"postgresql://root:123@localhost:5432/base?sslmode=disable"`
	LOG_LEVEL        string `env:"LOG_LEVEL" envDefault:"debug"`
	SECRET           string `env:"SECRET" envDefault:"your_secret_key"`
	GOOGLE_SECRET    string `env:"GOOGLE_SECRET" envDefault:"xxxxx-xxxxxx-xxxx"`
	GOOGLE_CLIENT_ID string `env:"GOOGLE_CLIENT_ID" envDefault:"xxxxx.apps.googleusercontent.com"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
