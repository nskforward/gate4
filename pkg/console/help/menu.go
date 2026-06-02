package help

import (
	"bytes"
	"io"
	"strings"

	"github.com/nskforward/gate4/pkg/console/output"
)

type Menu struct {
	title    string
	sections []*Section
}

type Section struct {
	name     string
	commands [][2]string
	maxKey   int
}

func NewMenu(title string) *Menu {
	return &Menu{
		title:    title,
		sections: make([]*Section, 0, 8),
	}
}

func newSection(name string) *Section {
	return &Section{
		name:     name,
		commands: make([][2]string, 0, 16),
	}
}

func (m *Menu) AddSection(name string) *Section {
	s := newSection(name)
	m.sections = append(m.sections, s)
	return s
}

func (s *Section) AddCommand(command, description string) {
	s.commands = append(s.commands, [2]string{command, description})
	if len(command) > s.maxKey {
		s.maxKey = len(command)
	}
}

func (m *Menu) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	buf.WriteByte('\n')
	buf.WriteString(output.BuildPrefix(output.Cyan, output.Bold))
	buf.WriteString(m.title)
	buf.WriteString(output.Reset)
	buf.WriteByte('\n')

	for _, section := range m.sections {
		buf.WriteByte('\n')
		buf.WriteString(output.BuildPrefix(output.Gray100))
		buf.WriteString(strings.ToTitle(section.name))
		buf.WriteByte(' ')
		buf.WriteString("COMMANDS")
		buf.WriteString(output.Reset)
		buf.WriteString("\n\n")
		for _, command := range section.commands {
			buf.WriteString("    ")
			buf.WriteString(output.BuildPrefix(output.White))
			buf.WriteString(command[0])
			buf.WriteString(output.Reset)
			buf.WriteByte(' ')
			buf.WriteString(strings.Repeat(" ", section.maxKey-len(command[0])))
			buf.WriteString(output.BuildPrefix(output.Gray100))
			buf.WriteString("- ")
			buf.WriteString(command[1])
			buf.WriteString(output.Reset)
			buf.WriteByte('\n')
		}
	}
	buf.WriteByte('\n')
	buf.WriteByte('\n')

	return buf.WriteTo(w)
}
