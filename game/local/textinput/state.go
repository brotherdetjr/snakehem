package textinput

import (
	"errors"
	"snakehem/input/controller"
	"strings"
)

var AvailableChars = []rune{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	' ',
}

type TextInput struct {
	value      string
	label      string
	charIndex  int
	maxLength  int
	controller controller.Controller
	callback   func(string)
	validation func(string) error
	error      error
}

func NewTextInput(value string, label string, maxLength int, controller controller.Controller, callback func(string)) *TextInput {
	return &TextInput{
		value:      value,
		label:      label,
		charIndex:  0,
		maxLength:  maxLength,
		controller: controller,
		callback:   callback,
		validation: nil,
		error:      nil,
	}
}

func (t *TextInput) ValidateNotEmpty(msg string) *TextInput {
	t.validation = func(text string) error {
		if strings.TrimSpace(t.value) == "" {
			return errors.New(msg)
		}
		return nil
	}
	return t
}
