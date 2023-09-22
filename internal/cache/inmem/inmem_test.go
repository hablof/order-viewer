package inmem

import (
	"strconv"
	"sync"
	"testing"

	"github.com/hablof/order-viewer/internal/models"
	"github.com/hablof/order-viewer/internal/service"
)

type inMemCacheRW struct {
	mu sync.RWMutex

	// map [OrderUid] -> Order
	kv map[string]models.Order
}

func (c *inMemCacheRW) Get(OrderId string) (models.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	order, ok := c.kv[OrderId]
	if !ok {
		return models.Order{}, service.ErrNotFound
	}

	return order, nil
}

func (c *inMemCacheRW) Set(order models.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.kv[order.OrderUID]; ok {
		return service.ErrDuplicatesNotAllowed
	}

	c.kv[order.OrderUID] = order
	return nil
}

func Benchmark_StandartMutex(b *testing.B) {

	c := inMemCache{
		mu: sync.Mutex{},
		kv: map[string]models.Order{},
	}

	testOrders := make([]models.Order, 10000)
	for i := range testOrders {
		s := strconv.Itoa(i)
		testOrders[i].OrderUID = s
	}

	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}

		for i, o := range testOrders {
			wg.Add(2)

			go func(o models.Order) {
				if err := c.Set(o); err != nil {
					b.Fatal(err)
				}
				wg.Done()
			}(o)
			go func(i int) {
				_, _ = c.Get(strconv.Itoa(i))
				wg.Done()
			}(i)
		}
		wg.Wait()
		c.mu.Lock()
		c.kv = map[string]models.Order{}
		c.mu.Unlock()
	}
}

func Benchmark_RWMutex(b *testing.B) {

	c := inMemCacheRW{
		mu: sync.RWMutex{},
		kv: map[string]models.Order{},
	}

	testOrders := make([]models.Order, 10000)
	for i := range testOrders {
		s := strconv.Itoa(i)
		testOrders[i].OrderUID = s
	}

	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}

		for i, o := range testOrders {
			wg.Add(2)

			go func(o models.Order) {
				if err := c.Set(o); err != nil {
					b.Fatal(err)
				}
				wg.Done()
			}(o)
			go func(i int) {
				_, _ = c.Get(strconv.Itoa(i))
				wg.Done()
			}(i)
		}
		wg.Wait()

		c.mu.Lock()
		c.kv = map[string]models.Order{}
		c.mu.Unlock()
	}
}
