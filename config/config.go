package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Postgres Postgres
		Nats     Nats
	}

	Postgres struct {
		URL string `env-required:"true" env:"POSTGRES_URL"`
	}

	Nats struct {
		Cluster string `env-required:"true" env:"NATS_CLUSTER"`
		URL     string `env-required:"true" env:"NATS_URL"`
	}
)

func NewConfig() (Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("config: %w", err)
	}
	return cfg, nil
}
