package textinput

import (
	"errors"
	"snakehem/input/controller"
	"strings"
)

const (
	KeyboardCols = 10
	KeyboardRows = 4
)

type SpecialKey int

const (
	SpecialKeyNone SpecialKey = iota
	SpecialKeyOK
	SpecialKeyDEL
)

type KeyboardKey struct {
	char       rune
	special    SpecialKey
	displayStr string
}

var keyboardGrid [KeyboardRows][KeyboardCols]*KeyboardKey

var AvailableChars = []rune{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	' ',
}

type TextInput struct {
	value      string
	label      string
	cursorRow  int
	cursorCol  int
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
		cursorRow:  0,
		cursorCol:  0,
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

func init() {
	initKeyboardGrid()
}

func initKeyboardGrid() {
	// Map regular characters from AvailableChars to grid
	for i, char := range AvailableChars {
		row := i / KeyboardCols
		col := i % KeyboardCols

		if row < KeyboardRows && col < KeyboardCols {
			displayStr := string(char)
			if char == ' ' {
				displayStr = "SPACE"
			}
			keyboardGrid[row][col] = &KeyboardKey{
				char:       char,
				special:    SpecialKeyNone,
				displayStr: displayStr,
			}
		}
	}

	// Add special keys
	// OK button at Row 3, Col 7
	keyboardGrid[3][9] = &KeyboardKey{
		char:       0,
		special:    SpecialKeyOK,
		displayStr: "OK",
	}

	// DEL button at Row 3, Col 8
	keyboardGrid[3][7] = &KeyboardKey{
		char:       0,
		special:    SpecialKeyDEL,
		displayStr: "DEL",
	}

	// Row 3, Col 8 remains nil (empty)
}

func (t *TextInput) getCurrentKey() *KeyboardKey {
	return keyboardGrid[t.cursorRow][t.cursorCol]
}

func isValidPosition(row, col int) bool {
	if row < 0 || row >= KeyboardRows || col < 0 || col >= KeyboardCols {
		return false
	}
	return keyboardGrid[row][col] != nil
}
