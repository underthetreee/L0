package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/underthetreee/L0/config"
	"github.com/underthetreee/L0/internal/model"
	"github.com/underthetreee/L0/internal/repository"
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

	sub, err := nats.NewSubcriber(cfg.Nats.Cluster, "orders-sub", cfg.Nats.URL)
	if err != nil {
		return err
	}
	defer sub.Close()

	repo := repository.NewPostgresRepository(db)

	msgCh, err := sub.Subscribe(ctx, "orders")
	if err != nil {
		return err
	}

	var order model.Order
	for msg := range msgCh {
		if err := json.Unmarshal(msg, &order); err != nil {
			log.Println("invalid json:", err)
			continue
		}

		log.Println("receiving order...")
		if err := repo.Store(ctx, order); err != nil {
			return err
		}
	}
	return nil
}
