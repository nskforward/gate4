package di

import "reflect"

func Resolve[T any](c *Container) (T, bool) {
	targetType := reflect.TypeFor[T]()
	val, ok := resolveAny(c, targetType)
	if !ok {
		var zero T
		return zero, false
	}
	return val.(T), true
}

func resolveAny(c *Container, t reflect.Type) (any, bool) {
	c.mu.RLock()
	// Проверяем кэш
	if val, ok := c.cache[t]; ok {
		c.mu.RUnlock()
		return val, true
	}
	prov, ok := c.providers[t]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}

	fnType := prov.fn.Type()
	args := make([]reflect.Value, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		argType := fnType.In(i)
		dep, ok := resolveAny(c, argType)
		if !ok {
			return nil, false
		}
		args[i] = reflect.ValueOf(dep)
	}

	// Вызываем конструктор
	results := prov.fn.Call(args)
	var instance any
	if prov.hasError {
		if !results[1].IsNil() {
			return nil, false
		}
		instance = results[0].Interface()
	} else {
		instance = results[0].Interface()
	}

	c.mu.Lock()
	c.cache[t] = instance
	c.mu.Unlock()
	return instance, true
}
