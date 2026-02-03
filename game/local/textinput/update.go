package textinput

import (
	"math"
	"snakehem/util"
	"unicode"
)

func (t *TextInput) Update() {
	c := t.controller
	if c.IsAnyJustPressed() {
		t.error = nil
	}
	if c.IsUpPressed() {
		t.moveUp()
	} else if c.IsDownPressed() {
		t.moveDown()
	} else if c.IsLeftPressed() {
		t.moveLeft()
	} else if c.IsRightPressed() {
		t.moveRight()
	} else if k := t.getCurrentKey(); k != nil && k.special == SpecialKeyEnter && c.IsStartJustPressed() {
		t.Submit()
	} else if c.IsStartPressed() {
		t.pressCurrentKey()
	}
}

func (t *TextInput) Submit() {
	if err := t.validation(t.value); err != nil {
		t.error = err
	} else {
		t.callback(t.value)
	}
}

func (t *TextInput) Clear() {
	t.value = ""
	if t.nameCapitalisation {
		t.WithCapsMode(true)
	}
}

func (t *TextInput) DeleteLastChar() {
	if t.value != "" {
		t.value = t.value[:len(t.value)-1]
	}
	if t.nameCapitalisation {
		t.WithCapsMode(t.value == "" || t.value[len(t.value)-1] == ' ' || t.value[len(t.value)-1] == '-')
	}
}

func (t *TextInput) AddSelectedChar() {
	key := t.getCurrentKey()
	if key != nil && key.char != 0 && len(t.value) < t.maxLength {
		t.value += string(key.char)
		if t.nameCapitalisation {
			if key.special == SpecialKeySpace || key.char == '-' {
				t.WithCapsMode(true)
			} else if unicode.IsUpper(key.char) {
				t.WithCapsMode(false)
			}
		}
	}
}

func (t *TextInput) ToggleCapsMode() {
	t.capsMode = !t.capsMode
	t.initKeyboardGrid()
}

func (t *TextInput) GetCapsMode() bool {
	return t.capsMode
}

func (t *TextInput) moveLeft() {
	for {
		t.cursorCol--
		if t.cursorCol < 0 {
			t.cursorCol = t.keyboardCols - 1
		}
		if t.keyboardGrid[t.cursorRow][t.cursorCol] != nil {
			break
		}
	}
}

func (t *TextInput) moveRight() {
	for {
		t.cursorCol++
		if t.cursorCol >= t.keyboardCols {
			t.cursorCol = 0
		}
		if t.keyboardGrid[t.cursorRow][t.cursorCol] != nil {
			break
		}
	}
}

func (t *TextInput) moveUp() {
	for {
		t.cursorRow--
		if t.cursorRow < 0 {
			t.cursorRow = t.keyboardRows - 1
		}
		if nearestCol := t.findNearestCol(); nearestCol != -1 {
			t.cursorCol = nearestCol
			break
		}
	}
}

func (t *TextInput) moveDown() {
	for {
		t.cursorRow++
		if t.cursorRow >= t.keyboardRows {
			t.cursorRow = 0
		}
		if nearestCol := t.findNearestCol(); nearestCol != -1 {
			t.cursorCol = nearestCol
			break
		}
	}
}

func (t *TextInput) findNearestCol() int {
	nearestCol := -1
	colDistance := math.MaxInt
	for col := 0; col < t.keyboardCols; col++ {
		if t.keyboardGrid[t.cursorRow][col] != nil {
			distance := util.AbsInt(col - t.cursorCol)
			if distance < colDistance {
				nearestCol = col
				colDistance = distance
			} else {
				break
			}
		}
	}
	return nearestCol
}

func (t *TextInput) pressCurrentKey() {
	key := t.getCurrentKey()
	if key == nil {
		return
	}

	switch key.special {
	case SpecialKeyEnter:
		// SpecialKeyEnter is handled specially to
		// prevent accidental submission
	case SpecialKeyClear:
		t.Clear()
	case SpecialKeyDel:
		t.DeleteLastChar()
	case SpecialKeyCaps:
		t.ToggleCapsMode()
	default:
		// Regular character
		t.AddSelectedChar()
	}
}

func (t *TextInput) getCurrentKey() *KeyboardKey {
	return t.keyboardGrid[t.cursorRow][t.cursorCol]
}
