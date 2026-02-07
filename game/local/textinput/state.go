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

type CapsBehaviour int

const (
	CapsBehaviourNormal CapsBehaviour = iota
	CapsBehaviourDisable
	CapsBehaviourNames
)

type keyPos struct {
	row, col int
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
	textColour     color.Color
	capsMode       bool
	capsBehaviour  CapsBehaviour
	cursorShown    bool
	cursorBlinkHz  float64
	spaceKeyPos    *keyPos
	delKeyPos      *keyPos
	enterKeyPos    *keyPos
	capsKeyPos     *keyPos
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
		keyboardCols:   getMinKeyCols(true, CapsBehaviourNormal),
		textColour:     color.White,
		capsMode:       false,
		capsBehaviour:  CapsBehaviourNormal,
		cursorShown:    false,
		cursorBlinkHz:  2,
		spaceKeyPos:    nil,
		delKeyPos:      nil,
		enterKeyPos:    nil,
		capsKeyPos:     nil,
	}
	t.initKeyboardGrid()
	return t
}

func (t *TextInput) WithValue(value string) *TextInput {
	t.value = value
	t.updateCaps()
	return t
}

func (t *TextInput) WithLabel(label string) *TextInput {
	t.label = label
	return t
}

func (t *TextInput) WithMaxLength(maxLength int) *TextInput {
	if maxLength < 0 || maxLength > 26 {
		panic("invalid max length for TextInput. Must be between 0 and 26")
	}
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

func (t *TextInput) WithCapsBehaviour(capsBehaviour CapsBehaviour) *TextInput {
	t.capsBehaviour = capsBehaviour
	t.updateCaps()
	return t
}

func (t *TextInput) WithCursorBlinkHz(cursorBlinkHz float64) *TextInput {
	t.cursorBlinkHz = cursorBlinkHz
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

func (t *TextInput) GetCurrentKey() *KeyboardKey {
	return t.keyboardGrid[t.cursorRow][t.cursorCol]
}

func (t *TextInput) Submit() {
	t.error = nil
	if err := t.validation(t.value); err != nil {
		t.error = err
	} else {
		t.callback(t.value)
	}
}

func (t *TextInput) Clear() {
	t.error = nil
	t.value = ""
	t.updateCaps()
}

func (t *TextInput) DeleteLastChar() {
	if t.value != "" {
		t.error = nil
		t.value = t.value[:len(t.value)-1]
		t.updateCaps()
	}
}

func (t *TextInput) AddSelectedChar() {
	key := t.GetCurrentKey()
	if key != nil && key.char != 0 && len(t.value) < t.maxLength {
		t.error = nil
		t.value += string(key.char)
		t.updateCaps()
	}
}

func (t *TextInput) ToggleCapsMode() {
	t.error = nil
	t.capsMode = !t.capsMode
	t.initKeyboardGrid()
}

func (t *TextInput) GetCapsMode() bool {
	return t.capsMode
}

func (t *TextInput) initKeyboardGrid() {
	if t.keyboardCols < getMinKeyCols(t.spaceAvailable, t.capsBehaviour) {
		panic("not enough keyboard cols")
	}

	// Map regular characters from AvailableChars to grid
	t.keyboardRows = int(math.Ceil(float64(len(t.availableChars))/float64(t.keyboardCols)) + 1)
	t.keyboardGrid = make([][]*KeyboardKey, t.keyboardRows)
	for i := range t.keyboardGrid {
		t.keyboardGrid[i] = make([]*KeyboardKey, t.keyboardCols)
	}

	for i, char := range t.availableChars {
		if unicode.IsLetter(char) {
			if t.capsMode {
				char = unicode.ToUpper(char)
			} else { // in case if available chars are listed in upper case
				char = unicode.ToLower(char)
			}
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

	col := 1
	// Add special keys
	t.keyboardGrid[0][col] = &KeyboardKey{
		char:       0,
		special:    SpecialKeyEnter,
		displayStr: "ENTER",
	}
	t.enterKeyPos = &keyPos{0, col}
	col = col + 2

	t.keyboardGrid[0][col] = &KeyboardKey{
		char:       0,
		special:    SpecialKeyClear,
		displayStr: "CLEAR",
	}
	col = col + 2

	t.keyboardGrid[0][col] = &KeyboardKey{
		char:       0,
		special:    SpecialKeyDel,
		displayStr: "DEL",
	}
	t.delKeyPos = &keyPos{0, col}
	col = col + 2

	if t.capsBehaviour != CapsBehaviourDisable {
		t.keyboardGrid[0][col] = &KeyboardKey{
			char:       0,
			special:    SpecialKeyCaps,
			displayStr: "CAPS",
		}
		t.capsKeyPos = &keyPos{0, col}
		col = col + 2
	} else {
		t.capsKeyPos = nil
	}

	if t.spaceAvailable {
		t.keyboardGrid[0][col] = &KeyboardKey{
			char:       ' ',
			special:    SpecialKeySpace,
			displayStr: "SPACE",
		}
		t.spaceKeyPos = &keyPos{0, col}
	} else {
		t.spaceKeyPos = nil
	}
}

func getMinKeyCols(spaceAvailable bool, capsBehaviour CapsBehaviour) int {
	// Special keys go in the first row and are separated by
	// empty spaces (nil keys) from each other.
	// Outermost keys are not placed at the edge.
	minCols := int(specialKeyCount)*2 + 1
	if capsBehaviour == CapsBehaviourDisable {
		minCols = minCols - 2
	}
	if !spaceAvailable {
		minCols = minCols - 2
	}
	return minCols
}
