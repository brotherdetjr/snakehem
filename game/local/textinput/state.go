package textinput

import (
	"errors"
	"math"
	"snakehem/input/controller"
	"strings"
)

type SpecialKey int

const (
	SpecialKeyNone SpecialKey = iota
	SpecialKeyEnter
	SpecialKeyClear
	SpecialKeyDel
	SpecialKeySpace
)

type KeyboardKey struct {
	char       rune
	special    SpecialKey
	displayStr string
}

var AZ09 = []rune{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
	'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1', '2', '3',
	'4', '5', '6', '7', '8', '9',
}

type TextInput struct {
	value          string
	label          string
	cursorRow      int
	cursorCol      int
	maxLength      int
	controller     controller.Controller
	callback       func(string)
	validation     func(string) error
	error          error
	availableChars []rune
	keyboardGrid   [][]*KeyboardKey
	spaceAvailable bool
	keyboardCols   int
	keyboardRows   int
}

func NewTextInput(controller controller.Controller) *TextInput {
	t := &TextInput{
		value:          "",
		label:          "",
		cursorRow:      0,
		cursorCol:      1, // row 0 col 1 -> SPACE key
		maxLength:      24,
		controller:     controller,
		callback:       func(string) {},
		validation:     nil,
		error:          nil,
		availableChars: AZ09,
		spaceAvailable: true,
		keyboardCols:   10,
	}
	t.initKeyboardGrid()
	return t
}

func (t *TextInput) WithValue(value string) *TextInput {
	t.value = value
	return t
}

func (t *TextInput) WithLabel(label string) *TextInput {
	t.label = label
	return t
}

func (t *TextInput) WithMaxLength(maxLength int) *TextInput {
	t.maxLength = maxLength
	return t
}

func (t *TextInput) WithCallback(callback func(string)) *TextInput {
	t.callback = callback
	return t
}

func (t *TextInput) WithAvailableChars(availableChars []rune) *TextInput {
	t.availableChars = availableChars
	t.initKeyboardGrid()
	return t
}

func (t *TextInput) WithSpaceAvailable(spaceAvailable bool) *TextInput {
	t.spaceAvailable = spaceAvailable
	t.initKeyboardGrid()
	return t
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

func (t *TextInput) initKeyboardGrid() {
	// Map regular characters from AvailableChars to grid
	t.keyboardRows = int(math.Ceil(float64(len(t.availableChars))/float64(t.keyboardCols)) + 1)
	t.keyboardGrid = make([][]*KeyboardKey, t.keyboardRows)
	for i := range t.keyboardGrid {
		t.keyboardGrid[i] = make([]*KeyboardKey, t.keyboardCols)
	}

	for i, char := range t.availableChars {
		row := i/t.keyboardCols + 1
		col := i % t.keyboardCols
		if row < t.keyboardRows && col < t.keyboardCols {
			t.keyboardGrid[row][col] = &KeyboardKey{
				char:       char,
				special:    SpecialKeyNone,
				displayStr: string(char),
			}
		}
	}

	// Add special keys
	t.keyboardGrid[0][1] = &KeyboardKey{
		char:       0,
		special:    SpecialKeyEnter,
		displayStr: "ENTER",
	}

	t.keyboardGrid[0][3] = &KeyboardKey{
		char:       0,
		special:    SpecialKeyClear,
		displayStr: "CLEAR",
	}

	t.keyboardGrid[0][5] = &KeyboardKey{
		char:       0,
		special:    SpecialKeyDel,
		displayStr: "DEL",
	}

	if t.spaceAvailable {
		t.keyboardGrid[0][7] = &KeyboardKey{
			char:       ' ',
			special:    SpecialKeySpace,
			displayStr: "SPACE",
		}
	}
}

func (t *TextInput) getCurrentKey() *KeyboardKey {
	return t.keyboardGrid[t.cursorRow][t.cursorCol]
}

func (t *TextInput) isValidPosition(row, col int) bool {
	if row < 0 || row >= t.keyboardRows || col < 0 || col >= t.keyboardCols {
		return false
	}
	return t.keyboardGrid[row][col] != nil
}
