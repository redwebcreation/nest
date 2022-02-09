package util

import (
	"fmt"
	"os"
)

// PrintErr prints the error message but does not exit, use PrintErrE for exit
func PrintErr(a ...interface{}) {
	if len(a) == 0 {
		_, _ = fmt.Fprintln(os.Stderr)
		return
	}
	_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf(a[0].(string), a[1:]...))
}

// PrintErrE prints the error message to stderr and exits the program.
func PrintErrE(a ...interface{}) {
	PrintErr(a...)

	os.Exit(1)
}
