package cache

import (
	"time"

	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/go-redis/redis/v7"
)

type redisCache struct {
	client    *redis.Client
	namespace string
}

func NewRedis(ns string) Cache {
	config := config.Get()
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisURL,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})
	return &redisCache{
		client:    client,
		namespace: ns,
	}
}

func (r *redisCache) Get(k string) interface{} {
	k = applyNamespace(r.namespace, k)
	v, err := r.client.Get(k).Result()
	if err != nil {
		return nil
	}
	return v
}

func (r *redisCache) Set(k string, v interface{}, d time.Duration) {
	k = applyNamespace(r.namespace, k)
	r.client.Set(k, v, d)
}

func (r *redisCache) Delete(k string) {
	k = applyNamespace(r.namespace, k)
	r.client.Del(k)
}
