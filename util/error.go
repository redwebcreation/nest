package util

import (
	"fmt"
	"os"
)

// PrintE prints the error message but does not exit, use FatalE for exit
func PrintE(a ...any) {
	if len(a) == 0 {
		_, _ = fmt.Fprintln(os.Stderr)
		return
	}

	if _, ok := a[0].(string); ok {
		_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf(a[0].(string), a[1:]...))
	} else {
		_, _ = fmt.Fprintln(os.Stderr, a...)
	}
}

// FatalE prints the error message to stderr and exits the program.
func FatalE(a ...any) {
	PrintE(a...)

	os.Exit(1)
}
