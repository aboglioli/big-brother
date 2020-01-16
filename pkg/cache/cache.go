package cache

import (
	"time"

	"github.com/aboglioli/big-brother/pkg/errors"
)

var (
	ErrCacheNotFound = errors.Internal.New("cache.not_found")
	ErrCacheSet      = errors.Internal.New("cache.set")
	ErrCacheDelete   = errors.Internal.New("cache.delete")
)

type Cache interface {
	Get(k string) (interface{}, error)
	Set(k string, v interface{}, d time.Duration) error
	Delete(k string) error
}
