package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/stan.go"
)

type NatsSubscriber struct {
	conn  stan.Conn
	msgCh chan []byte
}

func NewSubcriber(clusterID string, clientID string, url string) (*NatsSubscriber, error) {
	conn, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	sub := &NatsSubscriber{
		conn:  conn,
		msgCh: make(chan []byte),
	}
	return sub, nil
}

func (s *NatsSubscriber) Subscribe(ctx context.Context, subject string) error {
	if _, err := s.conn.Subscribe(subject, func(m *stan.Msg) {
		s.msgCh <- m.Data
	}, stan.DurableName("dur"), stan.DeliverAllAvailable()); err != nil {
		return err
	}
	return nil
}

func (s *NatsSubscriber) Messages() <-chan []byte {
	return s.msgCh
}

func (s *NatsSubscriber) Close() error {
	return s.conn.Close()
}
