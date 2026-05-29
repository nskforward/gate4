package di

import (
	"reflect"
	"sync"
)

type Container struct {
	mu        sync.RWMutex
	providers map[reflect.Type]*provider
	cache     map[reflect.Type]any
}

func NewContainer() *Container {
	return &Container{
		providers: make(map[reflect.Type]*provider),
		cache:     make(map[reflect.Type]any),
	}
}
