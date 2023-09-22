package inmem

import (
	"sync"

	"github.com/hablof/order-viewer/internal/models"
	"github.com/hablof/order-viewer/internal/service"
)

type inMemCache struct {
	// Из условия задачи не ясно соотношение
	// количества операций чтения и записи.
	//
	// Использую обычный Мьютекс, предполагая,
	// что количество чтений не превышает
	// количество записей многократно
	mu sync.Mutex

	// map [OrderUid] -> Order
	kv map[string]models.Order
}

func NewInMemCache() *inMemCache {
	c := &inMemCache{
		mu: sync.Mutex{},
		kv: map[string]models.Order{},
	}

	return c
}

func (c *inMemCache) LoadCache(cache map[string]models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.kv = cache
}

func (c *inMemCache) Get(OrderId string) (models.Order, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	order, ok := c.kv[OrderId]
	if !ok {
		return models.Order{}, service.ErrNotFound
	}

	return order, nil
}

func (c *inMemCache) Set(order models.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.kv[order.OrderUID]; ok {
		return service.ErrDuplicatesNotAllowed
	}

	c.kv[order.OrderUID] = order
	return nil
}
