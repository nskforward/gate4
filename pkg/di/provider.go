package di

import (
	"fmt"
	"reflect"
)

type provider struct {
	fn       reflect.Value
	hasError bool
}

func Provide[T any](c *Container, fn any) {
	targetType := reflect.TypeFor[T]()
	v := reflect.ValueOf(fn)
	t := v.Type()
	if t.Kind() != reflect.Func {
		panic(fmt.Sprintf("di: Provide requires func, got %v", t))
	}
	if t.NumOut() < 1 || t.NumOut() > 2 {
		panic(fmt.Sprintf("di: func must return 1 (T) or 2 (T, error), got %d result values", t.NumOut()))
	}
	if t.NumOut() == 2 && !t.Out(1).Implements(reflect.TypeFor[error]()) {
		panic("di: second result value must be an error")
	}

	implType := t.Out(0)

	if targetType.Kind() == reflect.Interface {
		if !implType.Implements(targetType) {
			panic(fmt.Sprintf("di: '%v' type does not implement '%v' interface", implType, targetType))
		}
	} else {
		if !implType.AssignableTo(targetType) {
			panic(fmt.Sprintf("di: '%v' type is not compatible with '%v' type", implType, targetType))
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.providers[targetType]; exists {
		panic(fmt.Sprintf("di: '%v' type is already registered", targetType))
	}
	c.providers[targetType] = &provider{
		fn:       v,
		hasError: t.NumOut() == 2,
	}
}

func ProvideValue[T any](c *Container, val T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	t := reflect.TypeFor[T]()
	if _, exists := c.providers[t]; exists {
		panic(fmt.Sprintf("di: '%v' type is already registered", t))
	}
	c.cache[t] = val
}
