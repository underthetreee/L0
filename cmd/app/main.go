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
	"github.com/underthetreee/L0/internal/service"
	"github.com/underthetreee/L0/pkg/cache"
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

	repo := repository.NewPostgresRepository(db)
	cch := cache.NewCache()

	log.Println("caching db...")
	if err := cch.LoadDB(ctx, repo); err != nil {
		return err
	}

	svc := service.NewOrderService(repo, cch)

	sub, err := nats.NewSubcriber(cfg.Nats.Cluster, "orders-sub", cfg.Nats.URL)
	if err != nil {
		return err
	}
	defer sub.Close()

	msgCh, err := sub.Subscribe(ctx, "orders")
	if err != nil {
		return err
	}

	var order model.Order
	for msg := range msgCh {
		log.Println("receiving order...")
		if err := json.Unmarshal(msg, &order); err != nil {
			log.Println("invalid json:", err)
			continue
		}

		if err := svc.Store(ctx, order); err != nil {
			return err
		}

	}
	return nil
}
