package cache

import (
	"time"

	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/go-redis/redis/v7"
)

var (
	ErrRedisConnect = errors.Internal.New("cache.redis.connect")
)

type redisCache struct {
	client    *redis.Client
	namespace string
}

func NewRedis(ns string) (*redisCache, error) {
	config := config.Get()
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisURL,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, ErrRedisConnect.Wrap(err)
	}

	return &redisCache{
		client:    client,
		namespace: ns,
	}, nil
}

func (r *redisCache) Get(k string) (interface{}, error) {
	k = applyNamespace(r.namespace, k)
	v, err := r.client.Get(k).Result()
	if err != nil {
		return nil, ErrCacheNotFound.M("key = %s", k).Wrap(err)
	}
	return v, nil
}

func (r *redisCache) Set(k string, v interface{}, d time.Duration) error {
	k = applyNamespace(r.namespace, k)
	res := r.client.Set(k, v, d)
	if res.Err() != nil {
		return ErrCacheSet.M("key = %s; value = %s", k, v).Wrap(res.Err())
	}
	return nil
}

func (r *redisCache) Delete(k string) error {
	k = applyNamespace(r.namespace, k)
	res := r.client.Del(k)
	if res.Err() != nil {
		return ErrCacheDelete.M("key = %s", k).Wrap(res.Err())
	}
	return nil
}

func (r *redisCache) Close() error {
	return r.client.Close()
}
