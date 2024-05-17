package main

import (
	"log"

	"github.com/nats-io/stan.go"
	"github.com/underthetreee/L0/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	pub, err := stan.Connect(cfg.Nats.Cluster, cfg.Nats.Client)
	if err != nil {
		return err
	}

	if err := pub.Publish("order", []byte("hello")); err != nil {
		return err
	}
	return nil
}
