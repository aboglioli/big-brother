package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

const (
	NoExpiration      = gocache.NoExpiration
	DefaultExpiration = gocache.DefaultExpiration
)

type goCache struct {
	cache     *gocache.Cache
	namespace string
}

func NewInMemory(ns string) Cache {
	c := gocache.New(2*time.Minute, 5*time.Minute)
	return &goCache{
		cache:     c,
		namespace: ns,
	}
}

func (c *goCache) Get(k string) interface{} {
	k = applyNamespace(c.namespace, k)
	data, ok := c.cache.Get(k)
	if !ok {
		return nil
	}
	return data
}

func (c *goCache) Set(k string, v interface{}, d time.Duration) {
	k = applyNamespace(c.namespace, k)
	c.cache.Set(k, v, d)
}

func (c *goCache) Delete(k string) {
	k = applyNamespace(c.namespace, k)
	c.cache.Delete(k)
}
