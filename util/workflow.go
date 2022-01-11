package util

import (
	"errors"
)

type Action func(...interface{}) ([]interface{}, error)

var (
	ErrSkipAction    = errors.New("skip action")
	ErrInvalidAction = errors.New("invalid action")
)

type Workflow []interface{}

func (w Workflow) Run(initialArgs ...interface{}) error {
	args := initialArgs
	var err error

	for k, action := range w {
		switch action := action.(type) {
		case Action:
			args, err = action(args...)
			if err == ErrSkipAction {
				continue
			}

			if err != nil {
				return err
			}
		case Workflow:
			err = action.Run(args...)
			if err != nil {
				return err
			}
		default:
			if len(w) == k+1 { // last action
				return ErrInvalidAction
			}

			args = append(args, action)
		}
	}

	return nil
}
