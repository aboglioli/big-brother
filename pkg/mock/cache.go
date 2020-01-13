package mock

import (
	"time"

	"github.com/aboglioli/big-brother/pkg/cache"
)

type Cache struct {
	Mock  Mock
	cache cache.Cache
}

func NewMockCache(ns string) *Cache {
	return &Cache{
		cache: cache.NewInMemory(ns),
	}
}

func (c *Cache) Get(k string) interface{} {
	call := Call("Get", k)
	v := c.cache.Get(k)
	c.Mock.Called(call.Return(v))
	return v
}

func (c *Cache) Set(k string, v interface{}, d time.Duration) {
	call := Call("Set", k, v, d)
	c.cache.Set(k, v, d)
	c.Mock.Called(call)
}

func (c *Cache) Delete(k string) {
	call := Call("Delete", k)
	c.cache.Delete(k)
	c.Mock.Called(call)
}
