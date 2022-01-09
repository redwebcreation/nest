package util

import "fmt"

func Prompt(prompt, defaultValue string, validator func(input string) bool) string {
	var input = defaultValue

	for {
		fmt.Print(prompt)
		if defaultValue != "" {
			fmt.Printf(" [%s]", defaultValue)
		}
		fmt.Print(": ")
		fmt.Scanln(&input)

		if validator == nil || validator(input) {
			break
		}
	}

	return input
}
