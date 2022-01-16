package util

import (
	"fmt"
	"io"
	"os"
)

// Stdin is a reader that reads from stdin.
var Stdin io.ReadWriter = os.Stdin

// Stdout is a writer that writes to stdout.
var Stdout io.ReadWriter = os.Stdout

// Prompt asks the user for input, validates it and returns it.
func Prompt(prompt, defaultValue string, validator func(input string) bool) string {
	var input = defaultValue

	for {
		_, _ = fmt.Fprint(Stdout, prompt)
		if defaultValue != "" {
			_, _ = fmt.Fprintf(Stdout, " [%s]", defaultValue)
		}
		_, _ = fmt.Fprint(Stdout, ": ")
		_, _ = fmt.Fscanln(Stdin, &input)

		if validator == nil || validator(input) {
			break
		}
	}

	return input
}
