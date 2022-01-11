package util

import (
	"fmt"
	"io"
	"os"
)

var Stdin io.Reader = os.Stdin

func Prompt(prompt, defaultValue string, validator func(input string) bool) string {
	var input = defaultValue

	for {
		fmt.Print(prompt)
		if defaultValue != "" {
			fmt.Printf(" [%s]", defaultValue)
		}
		fmt.Print(": ")
		fmt.Fscanln(Stdin, &input)

		if validator == nil || validator(input) {
			break
		}
	}

	return input
}
