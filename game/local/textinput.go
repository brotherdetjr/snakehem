package local

import (
	"image/color"
	"snakehem/assets/pxterm16"
	"snakehem/assets/pxterm24"
	"snakehem/game/common"
	"snakehem/input/controller"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
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
}

func NewTextInput(value string, label string, maxLength int, controller controller.Controller, callback func(string)) *TextInput {
	return &TextInput{
		value:      value,
		label:      label,
		charIndex:  0,
		maxLength:  maxLength,
		controller: controller,
		callback:   callback,
	}
}

func (t *TextInput) Submit() {
	t.callback(t.value)
}

func (t *TextInput) AddChar() {
	if len(t.value) == t.maxLength {
		return
	}
	t.value += string(AvailableChars[t.charIndex])
}

func (t *TextInput) DelChar() {
	if t.value == "" {
		return
	}
	t.value = t.value[:len(t.value)-1]
}

func (t *TextInput) NextChar() int {
	return (t.charIndex + 1) % len(AvailableChars)
}

func (t *TextInput) PrevChar() int {
	if t.charIndex == 0 {
		return len(AvailableChars) - 1
	}
	return t.charIndex - 1
}

func (t *TextInput) Update() {
	c := t.controller
	if c.IsUpJustPressed() {
		t.PrevChar()
	} else if c.IsDownJustPressed() {
		t.NextChar()
	} else if c.IsRightJustPressed() {
		t.AddChar()
	} else if c.IsLeftJustPressed() {
		t.DelChar()
	} else if c.IsStartJustPressed() {
		t.Submit()
	}
}

func (t *TextInput) Draw(screen *ebiten.Image) {
	// Draw semi-transparent overlay
	c := colornames.Darkolivegreen
	bg := color.RGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: 200,
	}
	screen.Fill(bg)

	// Draw title
	common.DrawTextCentered(
		screen,
		t.label,
		colornames.Yellow,
		common.GridDimPx/4,
		pxterm24.Font,
	)

	// Draw current name being entered
	currentNameY := common.GridDimPx / 2.5
	common.DrawTextCentered(
		screen,
		t.value,
		colornames.White,
		currentNameY,
		pxterm24.Font,
	)

	// Draw current character selection
	charY := currentNameY + float64(common.Pxterm24Height*2)
	currentChar := string(AvailableChars[t.charIndex])
	common.DrawTextCentered(
		screen,
		currentChar,
		colornames.Cyan,
		charY,
		pxterm24.Font,
	)

	// Draw instructions
	instructionsY := common.GridDimPx - float64(common.Pxterm16Height*6)
	common.DrawTextCentered(screen, "UP/DOWN: SELECT CHARACTER", colornames.Yellow, instructionsY, pxterm16.Font)
	common.DrawTextCentered(screen, "RIGHT: ADD CHARACTER", colornames.Yellow, instructionsY+float64(common.Pxterm16Height)*1.5, pxterm16.Font)
	common.DrawTextCentered(screen, "LEFT: DELETE CHARACTER", colornames.Yellow, instructionsY+float64(common.Pxterm16Height)*3, pxterm16.Font)
	common.DrawTextCentered(screen, "START: SUBMIT", colornames.Yellow, instructionsY+float64(common.Pxterm16Height)*4.5, pxterm16.Font)
}
