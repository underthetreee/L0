package cache

import (
	"context"
	"sync"

	"github.com/underthetreee/L0/internal/model"
)

type OrderGetter interface {
	GetAll(context.Context) ([]model.Order, error)
}

type Cache struct {
	data map[string]model.Order
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]model.Order),
	}
}

func (c *Cache) Set(orderID string, order model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[orderID] = order
}

func (c *Cache) Get(orderID string) (model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	order, ok := c.data[orderID]
	return order, ok
}

func (c *Cache) LoadDB(ctx context.Context, og OrderGetter) error {
	orders, err := og.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		c.Set(order.UID, order)
	}
	return nil
}
