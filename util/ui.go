package util

import (
	"fmt"
	"io"
	"os"
)

var Stdin io.ReadWriter = os.Stdin
var Stdout io.ReadWriter = os.Stdout

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
