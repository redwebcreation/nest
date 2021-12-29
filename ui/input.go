package ui

import "github.com/erikgeiser/promptkit/textinput"

type Input struct {
	Question   string
	Validation func(string) bool
	Default    string
}

func (i Input) Prompt() (string, error) {
	input := textinput.New(i.Question)
	input.Validate = i.Validation
	input.InitialValue = i.Default

	return input.RunPrompt()
}
