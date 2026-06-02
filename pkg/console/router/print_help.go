package router

import (
	"bytes"

	"github.com/nskforward/gate4/pkg/console/output"
)

func printHelpTitle(buf *bytes.Buffer, title string) {
	buf.WriteByte('\n')
	buf.WriteString(output.BuildPrefix(output.Cyan, output.Bold))
	buf.WriteString(title)
	buf.WriteString(output.Reset)
	buf.WriteByte('\n')
}

func printNodes(buf *bytes.Buffer, maxCmdLen int, nodes []*Node) {
	for _, node := range nodes {
		buf.WriteByte('\n')
		printCommands(buf, maxCmdLen, node)
	}
}

func printCommands(buf *bytes.Buffer, maxCmdLen int, node *Node) {
	if len(node.children) == 0 {
		if node.value != nil {
			buf.WriteString(output.BuildPrefix(output.White))
			buf.WriteString(node.value.Command)
			buf.WriteString(output.Reset)
			for range maxCmdLen - len(node.value.Command) {
				buf.WriteByte(' ')
			}
			buf.WriteString(" - ")
			buf.WriteString(output.BuildPrefix(output.Gray100))
			buf.WriteString(node.value.Desription)
			buf.WriteString(output.Reset)
			buf.WriteByte('\n')
		}
		return
	}

	for _, child := range node.children {
		printCommands(buf, maxCmdLen, child)
	}
}
