package util

import (
	"fmt"
	"os"
)

// PrintE prints the error message but does not exit, use FatalE for exit
func PrintE(a ...interface{}) {
	if len(a) == 0 {
		_, _ = fmt.Fprintln(os.Stderr)
		return
	}
	_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf(a[0].(string), a[1:]...))
}

// FatalE prints the error message to stderr and exits the program.
func FatalE(a ...interface{}) {
	PrintE(a...)

	os.Exit(1)
}
