package textinput

import (
	"math"
	"snakehem/game/common"
	"snakehem/model"
	"snakehem/util"
	"unicode"
)

func (t *TextInput) Update(ctx *common.Context) {
	t.cursorShown = ctx.Tick/int64(model.Tps/t.cursorBlinkHz)%2 == 0
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
	} else if k := t.GetCurrentKey(); k != nil && k.special == SpecialKeyEnter && c.IsStartJustPressed() {
		t.Submit()
	} else if c.IsStartPressed() {
		t.pressCurrentKey()
	}
}

func (t *TextInput) updateCaps() {
	if t.capsBehaviour == CapsBehaviourNames {
		t.WithCapsMode(t.value == "" || !unicode.IsLetter(rune(t.value[len(t.value)-1])))
	}
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
	key := t.GetCurrentKey()
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
