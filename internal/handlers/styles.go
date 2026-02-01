package handlers

import (
	"fmt"
	"strings"
)

const (
	Divider = "━━━━━━━━━━━━━━━━━━━━━━"
)

func ToBoldSans(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r - 'A' + 0x1D5D4)
		case r >= 'a' && r <= 'z':
			b.WriteRune(r - 'a' + 0x1D5EE)
		case r >= '0' && r <= '9':
			b.WriteRune(r - '0' + 0x1D7EC)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

type StyleBuilder struct {
	title  string
	icon   string
	rows   []string
	footer string
}

func (c *Context) NewStyle(title string, icon string) *StyleBuilder {
	return &StyleBuilder{
		title: title,
		icon:  icon,
	}
}

func (s *StyleBuilder) AddRow(key string, value interface{}) *StyleBuilder {
	return s.AddRowWithIcon("▫️", key, value)
}

func (s *StyleBuilder) AddRowWithIcon(icon string, key string, value interface{}) *StyleBuilder {
	s.rows = append(s.rows, fmt.Sprintf("%s %-12s : %v", icon, ToBoldSans(key), value))
	return s
}

func (s *StyleBuilder) SetFooter(footer string) *StyleBuilder {
	s.footer = footer
	return s
}

func (s *StyleBuilder) Build() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s %s\n", s.icon, ToBoldSans(strings.ToUpper(s.title))))
	b.WriteString(Divider + "\n")
	for _, row := range s.rows {
		b.WriteString(row + "\n")
	}
	b.WriteString(Divider)
	if s.footer != "" {
		b.WriteString("\n" + s.footer)
	}
	return b.String()
}
