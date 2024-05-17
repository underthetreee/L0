package nats

import (
	"fmt"
	"log"

	"github.com/nats-io/stan.go"
)

type NatsSubscriber struct {
	sub stan.Conn
}

func NewSubcriber(cluster string, client string) (NatsSubscriber, error) {
	sub, err := stan.Connect(cluster, client)
	if err != nil {
		return NatsSubscriber{}, fmt.Errorf("nats connect: %w", err)
	}
	return NatsSubscriber{
		sub: sub,
	}, nil
}

func (s *NatsSubscriber) Subscribe(subject string) error {
	sub, err := s.sub.Subscribe(subject, func(msg *stan.Msg) {
		log.Printf("message: %s\n", string(msg.Data))
	}, stan.DurableName("dur"))
	if err != nil {
		return fmt.Errorf("subscriber: %w", err)
	}
	defer sub.Close()
	return nil
}
