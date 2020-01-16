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

func (c *goCache) Get(k string) (interface{}, error) {
	k = applyNamespace(c.namespace, k)
	data, ok := c.cache.Get(k)
	if !ok {
		return nil, ErrCacheNotFound.C("key", k)
	}
	return data, nil
}

func (c *goCache) Set(k string, v interface{}, d time.Duration) error {
	k = applyNamespace(c.namespace, k)
	c.cache.Set(k, v, d)
	return nil
}

func (c *goCache) Delete(k string) error {
	k = applyNamespace(c.namespace, k)
	c.cache.Delete(k)
	return nil
}
