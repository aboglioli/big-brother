package cache

import "time"

type Cache interface {
	Get(k string) interface{}
	Set(k string, v interface{}, d time.Duration)
	Delete(k string)
}
