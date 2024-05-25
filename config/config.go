package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Postgres Postgres
		Nats     Nats
		HTTP     HTTP
	}

	Postgres struct {
		URL string `env-required:"true" env:"POSTGRES_URL"`
	}

	Nats struct {
		Cluster string `env-required:"true" env:"NATS_CLUSTER"`
		URL     string `env-required:"true" env:"NATS_URL"`
	}

	HTTP struct {
		Addr         string        `env-required:"true" env:"HTTP_ADDR"`
		ReadTimeout  time.Duration `env-required:"true" env:"HTTP_READ_TIMEOUT"`
		WriteTimeout time.Duration `env-required:"true" env:"HTTP_WRITE_TIMEOUT"`
	}
)

func NewConfig() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return &cfg, nil
}
