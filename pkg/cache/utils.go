package cache

import "fmt"

func applyNamespace(ns string, key string) string {
	if ns != "" {
		return fmt.Sprintf("%s:%s", ns, key)
	}
	return key
}
