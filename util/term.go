package util

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
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
