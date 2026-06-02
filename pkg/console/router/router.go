package router

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
)

type Router struct {
	title     string
	rootNode  *Node
	maxCmdLen int
}

func NewRouter(title string) *Router {
	return &Router{
		title:    title,
		rootNode: new(Node),
	}
}

func (r *Router) Handle(command, descriptiotn string, handler Handler) {
	segments := strings.Fields(command)
	if len(segments) == 0 {
		panic(fmt.Errorf("invalid CLI route command: %s", command))
	}

	if len(command) > r.maxCmdLen {
		r.maxCmdLen = len(command)
	}

	curr := r.rootNode
	for _, segment := range segments {
		if segment != "" && (segment[0] == '<' || segment[0] == '[') {
			break
		}
		curr = curr.CreateChild(segment)
	}

	curr.SetValue(&NodeValue{
		Desription: descriptiotn,
		Handler:    handler,
		Command:    command,
	})
}

func (r *Router) Run(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("unknown command")
	}

	args = os.Args[1:]

	curr := r.rootNode
	tail := args

	for _, arg := range args {
		next := curr.GetChild(arg)
		if next == nil {
			break
		}
		curr = next
		tail = tail[1:]
	}

	value, ok := curr.GetValue()
	if !ok || value.Handler == nil {
		return fmt.Errorf("unknown command")
	}

	return value.Handler(ctx, tail)
}

func (r *Router) PrintHelp() Handler {
	return func(ctx context.Context, args []string) error {
		var buf bytes.Buffer
		printHelpTitle(&buf, r.title)
		printNodes(&buf, r.maxCmdLen, r.rootNode.children)
		buf.WriteByte('\n')
		buf.WriteTo(os.Stdout)
		return nil
	}
}
