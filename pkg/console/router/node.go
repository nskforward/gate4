package router

import (
	"fmt"
)

type Node struct {
	key      string
	value    *NodeValue
	children []*Node
}

func (n *Node) GetChild(key string) *Node {
	for _, child := range n.children {
		if child.key == key {
			return child
		}
	}
	return nil
}

func (n *Node) SetValue(value *NodeValue) {
	if n.value != nil {
		panic(fmt.Errorf("cannot overwrite existing value in node with key: %s", n.key))
	}
	n.value = value
}

func (n *Node) GetValue() (NodeValue, bool) {
	if n.value == nil {
		return NodeValue{}, false
	}
	return *n.value, true
}

func (n *Node) CreateChild(key string) *Node {
	child := n.GetChild(key)
	if child != nil {
		return child
	}
	if n.children == nil {
		n.children = make([]*Node, 0, 4)
	}
	child = &Node{
		key: key,
	}
	n.children = append(n.children, child)
	return child
}
