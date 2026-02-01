package textinput

import (
	"image/color"
	"snakehem/assets/pxterm16"
	"snakehem/assets/pxterm24"
	"snakehem/game/common"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

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
	titleY := common.GridDimPx / 4.0
	common.DrawTextCentered(
		screen,
		t.label,
		colornames.Yellow,
		titleY,
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

	// Draw error if needed
	if t.error != nil {
		errorY := titleY + float64(common.Pxterm24Height+common.Pxterm16Height/2)
		common.DrawTextCentered(screen, strings.ToUpper(t.error.Error()), colornames.Orangered, errorY, pxterm16.Font)
	}

	// Draw instructions
	instructionsY := common.GridDimPx - float64(common.Pxterm16Height*6)
	common.DrawTextCentered(screen, "UP/DOWN: SELECT CHARACTER", colornames.Yellow, instructionsY, pxterm16.Font)
	common.DrawTextCentered(screen, "RIGHT: ADD CHARACTER", colornames.Yellow, instructionsY+float64(common.Pxterm16Height)*1.5, pxterm16.Font)
	common.DrawTextCentered(screen, "LEFT: DELETE CHARACTER", colornames.Yellow, instructionsY+float64(common.Pxterm16Height)*3, pxterm16.Font)
	common.DrawTextCentered(screen, "START: SUBMIT", colornames.Yellow, instructionsY+float64(common.Pxterm16Height)*4.5, pxterm16.Font)
}
