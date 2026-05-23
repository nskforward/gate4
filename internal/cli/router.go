package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nskforward/gate4/pkg/tree"
)

type Handler func(ctx context.Context, args []string) error

type Router struct {
	node *tree.Node[Handler]
}

func NewRouter() *Router {
	return &Router{
		node: new(tree.Node[Handler]),
	}
}

func (r *Router) Handle(path string, handler Handler) {
	segments := strings.Fields(path)
	if len(segments) == 0 {
		panic(fmt.Errorf("invalid CLI route path: %s", path))
	}
	curr := r.node
	for _, segment := range segments {
		curr = curr.CreateChild(segment)
	}
	curr.SetValue(handler)
}

func (r *Router) Run(ctx context.Context) error {
	if len(os.Args) < 2 {
		return fmt.Errorf("unknown command")
	}

	args := os.Args[1:]

	curr := r.node
	tail := args

	for _, arg := range args {
		next := curr.GetChild(arg)
		if next == nil {
			break
		}
		curr = next
		tail = tail[1:]
	}

	handler, ok := curr.GetValue()
	if !ok {
		return fmt.Errorf("unknown command")
	}
	return handler(ctx, tail)
}
