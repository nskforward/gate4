package di

import (
	"fmt"
	"reflect"
)

func Resolve[T any](c *Container) (T, error) {
	targetType := reflect.TypeFor[T]()
	val, err := resolveAny(c, targetType)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("di: resolving %v: %w", targetType, err)
	}
	return val.(T), nil
}

func resolveAny(c *Container, t reflect.Type) (any, error) {
	c.mu.RLock()
	if val, ok := c.cache[t]; ok {
		c.mu.RUnlock()
		return val, nil
	}
	prov, ok := c.providers[t]
	c.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no provider for type %v", t)
	}

	fnType := prov.fn.Type()
	args := make([]reflect.Value, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		argType := fnType.In(i)
		dep, err := resolveAny(c, argType)
		if err != nil {
			return nil, fmt.Errorf("argument %d (%v): %w", i, argType, err)
		}
		args[i] = reflect.ValueOf(dep)
	}

	results := prov.fn.Call(args)
	var instance any
	if prov.hasError {
		if !results[1].IsNil() {
			constructorErr := results[1].Interface().(error)
			return nil, constructorErr
		}
		instance = results[0].Interface()
	} else {
		instance = results[0].Interface()
	}

	c.mu.Lock()
	c.cache[t] = instance
	c.mu.Unlock()
	return instance, nil
}
