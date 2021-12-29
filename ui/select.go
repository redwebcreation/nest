package ui

import "github.com/erikgeiser/promptkit/selection"

type Select struct {
	Question string
	Choices  []string
}

func (s Select) Prompt() (string, error) {
	input := selection.New(s.Question, selection.Choices(s.Choices))
	input.PageSize = 5
	choice, err := input.RunPrompt()
	if err != nil {
		return "", err
	}

	return choice.String, nil
}
