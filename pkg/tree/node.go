package tree

import (
	"fmt"
)

type Node[T any] struct {
	key      string
	value    *T
	children []*Node[T]
}

func (n *Node[T]) GetChild(key string) *Node[T] {
	for _, child := range n.children {
		if child.key == key {
			return child
		}
	}
	return nil
}

func (n *Node[T]) SetValue(value T) {
	if n.value != nil {
		panic(fmt.Errorf("cannot overwrite existing value in node with key: %s", n.key))
	}
	n.value = &value
}

func (n *Node[T]) GetValue() (T, bool) {
	if n.value == nil {
		var def T
		return def, false
	}
	return *n.value, true
}

func (n *Node[T]) CreateChild(key string) *Node[T] {
	child := n.GetChild(key)
	if child != nil {
		return child
	}
	if n.children == nil {
		n.children = make([]*Node[T], 0, 4)
	}
	child = &Node[T]{
		key: key,
	}
	n.children = append(n.children, child)
	return child
}
