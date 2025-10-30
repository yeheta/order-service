package cache

import (
	"order-service/internal/domain"
	"sync"
)

type Cache interface {
	Set(uid string, order domain.Order)
	Get(uid string) (domain.Order, bool)
	GetAll() []domain.Order
	Size() int
}

type MemoryCache struct {
	data *sync.Map
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: &sync.Map{},
	}
}

func (c *MemoryCache) Set(uid string, order domain.Order) {
	c.data.Store(uid, order)
}

func (c *MemoryCache) Get(uid string) (domain.Order, bool) {
	val, ok := c.data.Load(uid)
	if !ok {
		return domain.Order{}, false
	}
	return val.(domain.Order), true
}

func (c *MemoryCache) GetAll() []domain.Order {
	var orders []domain.Order
	c.data.Range(func(key, value interface{}) bool {
		orders = append(orders, value.(domain.Order))
		return true
	})
	return orders
}

func (c *MemoryCache) Size() int {
	count := 0
	c.data.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}