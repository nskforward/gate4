package di

import (
	"fmt"
	"reflect"
)

func Resolve[T any](c *Container) T {
	targetType := reflect.TypeFor[T]()
	val := resolveAny(c, targetType)
	return val.(T)
}

func resolveAny(c *Container, t reflect.Type) any {
	c.mu.RLock()
	if val, ok := c.cache[t]; ok {
		c.mu.RUnlock()
		return val
	}
	prov, ok := c.providers[t]
	c.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("di: provider not found for '%v' type", t))
	}

	fnType := prov.fn.Type()
	args := make([]reflect.Value, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		argType := fnType.In(i)
		args[i] = reflect.ValueOf(resolveAny(c, argType))
	}

	results := prov.fn.Call(args)
	var instance any
	if prov.hasError {
		if !results[1].IsNil() {
			panic(fmt.Sprintf("di: constructor returned an error for '%v' type: %v", t, results[1].Interface()))
		}
		instance = results[0].Interface()
	} else {
		instance = results[0].Interface()
	}

	c.mu.Lock()
	c.cache[t] = instance
	c.mu.Unlock()
	return instance
}
