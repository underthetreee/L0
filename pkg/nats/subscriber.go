package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/stan.go"
)

type NatsSubscriber struct {
	conn stan.Conn
}

func NewSubcriber(clusterID string, clientID string, url string) (NatsSubscriber, error) {
	conn, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		return NatsSubscriber{}, fmt.Errorf("nats connect: %w", err)
	}
	return NatsSubscriber{
		conn: conn,
	}, nil
}

func (s *NatsSubscriber) Subscribe(ctx context.Context, subject string) (<-chan []byte, error) {
	msgCh := make(chan []byte)

	if _, err := s.conn.Subscribe(subject, func(m *stan.Msg) {
		msgCh <- m.Data
	}, stan.DurableName("dur"), stan.DeliverAllAvailable()); err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		close(msgCh)
	}()

	return msgCh, nil
}

func (s *NatsSubscriber) Close() error {
	return s.conn.Close()
}
