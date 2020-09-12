package cache

import (
	"errors"
)

type Cache interface {
	Get(key string) (interface{}, error)
	Put(key string, value interface{}) error
	Exists(key string) bool
	Invalidate(key string) error
}

type item struct {
	Object interface{}
}

var (
	ErrNotFound = errors.New("cache: key not found")
)
