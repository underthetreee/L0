package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/underthetreee/L0/config"
	"github.com/underthetreee/L0/pkg/nats"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	db, err := pgxpool.New(ctx, cfg.Postgres.URL)
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}
	defer db.Close()

	sub, err := nats.NewSubcriber(cfg.Nats.Cluster, cfg.Nats.Client)
	if err != nil {
		return err
	}
	if err := sub.Subscribe("order"); err != nil {
		return err
	}
	return nil
}
