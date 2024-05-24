package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/underthetreee/L0/config"
	"github.com/underthetreee/L0/internal/handler"
	"github.com/underthetreee/L0/internal/model"
	"github.com/underthetreee/L0/internal/repository"
	"github.com/underthetreee/L0/internal/service"

	"github.com/underthetreee/L0/pkg/cache"
	"github.com/underthetreee/L0/pkg/nats"
	"github.com/underthetreee/L0/pkg/server"
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

	if err := sub.Subscribe(ctx, "orders"); err != nil {
		return err
	}

	var (
		quitCh = make(chan os.Signal, 1)
		errCh  = make(chan error)
	)
	signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)

	orderHandler := handler.NewOrderHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/orders/{id}", orderHandler.Get)

	srv := server.NewServer(cfg, mux)

	go func() {
		if err := srv.Run(); err != nil {
			errCh <- err
		}
	}()

	for {
		select {
		case msg := <-sub.Messages():
			var order model.Order

			log.Println("receiving order...")
			if err := json.Unmarshal(msg, &order); err != nil {
				log.Println("invalid json:", err)
				continue
			}

			if err := svc.Store(ctx, order); err != nil {
				return err
			}
		case <-quitCh:
			if err := srv.Shutdown(ctx); err != nil {
				return err
			}
			return nil
		case err := <-errCh:
			return err
		}
	}
}
