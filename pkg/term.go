package pkg

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

var Red = lipgloss.NewStyle().
	Bold(true).
	Inline(true).
	Foreground(lipgloss.Color("#ff0000"))

var Yellow = lipgloss.NewStyle().
	Bold(true).
	Inline(true).
	Foreground(lipgloss.Color("#ffff00"))

var Green = lipgloss.NewStyle().
	Bold(true).
	Inline(true).
	Foreground(lipgloss.Color("#00ff00"))

var Gray = lipgloss.NewStyle().
	Bold(true).
	Inline(true).
	Foreground(lipgloss.Color("#808080"))

var White = lipgloss.NewStyle().
	Bold(true).
	Inline(true).
	Foreground(lipgloss.Color("#ffffff"))

var Stderr = output{
	file: os.Stderr,
}

type output struct {
	file *os.File
}

func (o output) Print(a ...any) {
	if len(a) == 0 {
		return
	}

	if _, format := a[0].(string); format {
		fmt.Fprintf(o.file, a[0].(string), a[1:]...)
	} else {
		fmt.Fprint(o.file, a...)
	}
}

func (o output) Println(a ...any) {
	fmt.Fprintln(o.file, a...)
}

func (o output) Fatal(format any, a ...any) {
	args := append([]any{format}, a...)

	o.Print(args...)

	os.Exit(1)
}

func Printf(style lipgloss.Style, format string, a ...any) {
	left := strings.TrimLeft(format, "\n")
	right := strings.TrimRight(format, "\n")

	leftCount := len(format) - len(left)
	rightCount := len(format) - len(right)

	if leftCount > 0 {
		fmt.Print(strings.Repeat("\n", leftCount))
	}

	fmt.Print(style.Render(fmt.Sprintf(right, a...)))

	if rightCount > 0 {
		fmt.Print(strings.Repeat("\n", rightCount))
	}
}
