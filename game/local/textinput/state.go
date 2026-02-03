package textinput

import (
	"errors"
	"image/color"
	"math"
	"snakehem/input/controller"
	"strings"
	"unicode"
)

type SpecialKey int

const (
	SpecialKeyNone SpecialKey = iota
	SpecialKeyEnter
	SpecialKeyClear
	SpecialKeyDel
	SpecialKeySpace
	SpecialKeyCaps
	specialKeyCount = iota - 1 // SpecialKeyNone doesn't count
)

type KeyboardKey struct {
	char       rune
	special    SpecialKey
	displayStr string
}

var AZ09 = []rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z', '0', '1', '2', '3',
	'4', '5', '6', '7', '8', '9',
}

type TextInput struct {
	value              string
	label              string
	cursorRow          int
	cursorCol          int
	maxLength          int
	controller         controller.Controller
	callback           func(string)
	validation         func(string) error
	error              error
	availableChars     []rune
	keyboardGrid       [][]*KeyboardKey
	spaceAvailable     bool
	keyboardCols       int
	keyboardRows       int
	textColour         color.Color
	capsMode           bool
	nameCapitalisation bool
}

func NewTextInput(controller controller.Controller) *TextInput {
	t := &TextInput{
		value:              "",
		label:              "",
		cursorRow:          0,
		cursorCol:          1, // row 0 col 1 -> SPACE key
		maxLength:          24,
		controller:         controller,
		callback:           func(string) {},
		validation:         nil,
		error:              nil,
		availableChars:     AZ09,
		spaceAvailable:     true,
		keyboardCols:       GetMinKeyCols(true),
		textColour:         color.White,
		capsMode:           false,
		nameCapitalisation: false,
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

func (t *TextInput) WithKeyboardCols(keyboardCols int) *TextInput {
	t.keyboardCols = keyboardCols
	t.initKeyboardGrid()
	return t
}

func (t *TextInput) WithTextColour(textColour color.Color) *TextInput {
	t.textColour = textColour
	return t
}

func (t *TextInput) WithCapsMode(capsMode bool) *TextInput {
	t.capsMode = capsMode
	t.initKeyboardGrid()
	return t
}

func (t *TextInput) WithNameCapitalisation(nameCapitalisation bool) *TextInput {
	t.nameCapitalisation = nameCapitalisation
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
	if t.keyboardCols < GetMinKeyCols(t.spaceAvailable) {
		panic("not enough keyboard cols")
	}
	// TODO validate we don't mix upper and lower case in the layout

	// Map regular characters from AvailableChars to grid
	t.keyboardRows = int(math.Ceil(float64(len(t.availableChars))/float64(t.keyboardCols)) + 1)
	t.keyboardGrid = make([][]*KeyboardKey, t.keyboardRows)
	for i := range t.keyboardGrid {
		t.keyboardGrid[i] = make([]*KeyboardKey, t.keyboardCols)
	}

	for i, char := range t.availableChars {
		if t.capsMode {
			char = unicode.ToUpper(char)
		}
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

	t.keyboardGrid[0][7] = &KeyboardKey{
		char:       0,
		special:    SpecialKeyCaps,
		displayStr: "CAPS",
	}

	if t.spaceAvailable {
		t.keyboardGrid[0][9] = &KeyboardKey{
			char:       ' ',
			special:    SpecialKeySpace,
			displayStr: "SPACE",
		}
	}
}

func GetMinKeyCols(spaceAvailable bool) int {
	// Special keys go in the first row and are separated by
	// empty spaces (nil keys) from each other.
	// Outermost keys are not placed at the edge.
	minCols := int(specialKeyCount)*2 + 1
	if !spaceAvailable {
		minCols = minCols - 2
	}
	return minCols
}
