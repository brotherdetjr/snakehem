package textinput

import (
	"math"
	"snakehem/game/common"
	"snakehem/input/controller"
	"snakehem/input/keyboard"
	"snakehem/input/keyboardwasd"
	"snakehem/model"
	"snakehem/util"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (t *TextInput) Update(ctx *common.Context) {
	t.cursorShown = ctx.Tick/int64(model.Tps/t.cursorBlinkHz)%2 == 0
	c := t.controller
	// For text input, replacing WASD controls with a "normal" keyboard,
	// because WASD are letters, but text input handles letters as direct input.
	if c == keyboardwasd.Instance {
		c = keyboard.Instance
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
	} else {
		// Direct keyboard input
		pressedKeys := inpututil.AppendPressedKeys(nil)
		for _, pressedKey := range pressedKeys {
			if !controller.IsRepeatingKeyboard(pressedKey) {
				continue
			}
			if len(pressedKey.String()) == 1 {
				// Try to find a matching key on the virtual keyboard
				pressedKeyName := []rune(pressedKey.String())[0]
				for r, gridRow := range t.keyboardGrid {
					for c, gridKey := range gridRow {
						if gridKey != nil && unicode.ToUpper(gridKey.char) == pressedKeyName {
							t.cursorRow = r
							t.cursorCol = c
							t.pressCurrentKey()
						}
					}
				}
			} else if pressedKey == ebiten.KeySpace && t.spaceAvailable {
				t.cursorRow = t.spaceKeyPos.row
				t.cursorCol = t.spaceKeyPos.col
				t.pressCurrentKey()
			} else if pressedKey == ebiten.KeyBackspace {
				t.cursorRow = t.delKeyPos.row
				t.cursorCol = t.delKeyPos.col
				t.DeleteLastChar()
			} else if pressedKey == ebiten.KeyEnter || pressedKey == ebiten.KeyNumpadEnter {
				t.cursorRow = t.enterKeyPos.row
				t.cursorCol = t.enterKeyPos.col
				t.Submit()
			}
		}
	}
}

func (t *TextInput) updateCaps() {
	if t.capsBehaviour == CapsBehaviourNames {
		t.WithCapsMode(t.value == "" || !unicode.IsLetter(rune(t.value[len(t.value)-1])))
	}
}

func (t *TextInput) moveLeft() {
	t.error = nil
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
	t.error = nil
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
	t.error = nil
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
	t.error = nil
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
