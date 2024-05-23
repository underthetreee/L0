package service

import (
	"context"

	"github.com/underthetreee/L0/internal/model"
)

type OrderStorer interface {
	Store(context.Context, model.Order) error
}

type Cache interface {
	Set(string, model.Order)
	Get(string) (model.Order, bool)
}

type OrderService struct {
	repo  OrderStorer
	cache Cache
}

func NewOrderService(repo OrderStorer, cache Cache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

func (s *OrderService) Store(ctx context.Context, order model.Order) error {
	if err := s.repo.Store(ctx, order); err != nil {
		return err
	}

	s.cache.Set(order.UID, order)
	return nil
}
