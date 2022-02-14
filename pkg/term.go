package pkg

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os"
)

var Red = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#ff0000"))

var Yellow = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#ffff00"))

var Green = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#00ff00"))

var Gray = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#808080"))

var White = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#ffffff"))

var Stderr = output{
	file: os.Stderr,
}

type output struct {
	file *os.File
}

func (o output) Print(format any, a ...any) {
	if _, ok := format.(string); ok {
		fmt.Fprintf(o.file, format.(string), a...)
	} else {
		fmt.Fprint(o.file, format)
		fmt.Fprint(o.file, a...)
	}
}

func (o output) Println(a ...any) {
	fmt.Fprintln(o.file, a...)
}

func (o output) Fatal(format any, a ...any) {
	o.Print(format, a...)

	os.Exit(1)
}
