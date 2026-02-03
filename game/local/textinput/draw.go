package textinput

import (
	"image/color"
	"snakehem/assets/pxterm16"
	"snakehem/assets/pxterm24"
	"snakehem/game/common"
	"snakehem/util"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func (t *TextInput) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Darkolivegreen)

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
	currentNameY := common.GridDimPx / 2.7
	common.DrawTextCentered(
		screen,
		"["+util.PadRight(t.value, t.maxLength)+"]",
		colornames.White,
		currentNameY,
		pxterm24.Font,
	)
	common.DrawTextCentered(
		screen,
		"["+util.PadRight("", t.maxLength)+"]",
		colornames.Orange,
		currentNameY,
		pxterm24.Font,
	)

	// Draw virtual keyboard grid
	keyboardStartY := currentNameY + float64(common.Pxterm24Height*3)
	t.drawKeyboardGrid(screen, keyboardStartY)

	// Draw error if needed
	if t.error != nil {
		errorY := titleY + float64(common.Pxterm24Height+common.Pxterm16Height/2)
		common.DrawTextCentered(screen, strings.ToUpper(t.error.Error()), colornames.Orangered, errorY, pxterm16.Font)
	}

	// Draw instructions
	instructionsY := common.GridDimPx - float64(common.Pxterm16Height)*2.5
	common.DrawTextCentered(screen, "ARROWS: NAVIGATE KEYBOARD", colornames.Yellow, instructionsY, pxterm16.Font)
	common.DrawTextCentered(screen, "START: PRESS SELECTED KEY", colornames.Yellow, instructionsY+float64(common.Pxterm16Height), pxterm16.Font)
}

func (t *TextInput) drawKeyboardGrid(screen *ebiten.Image, startY float64) {
	const keySpacingX = 70 // Horizontal spacing between keys
	const keySpacingY = 32 // Vertical spacing between rows

	// Calculate grid dimensions
	totalGridWidth := (t.keyboardCols - 1) * keySpacingX
	gridStartX := (common.GridDimPx - totalGridWidth) / 2

	for row := 0; row < t.keyboardRows; row++ {
		for col := 0; col < t.keyboardCols; col++ {
			key := t.keyboardGrid[row][col]
			if key == nil {
				continue // Skip empty cells
			}

			// Calculate position
			x := gridStartX + (col * keySpacingX)
			y := int(startY) + (row * keySpacingY)

			// Determine if this key is selected
			isSelected := row == t.cursorRow && col == t.cursorCol

			// Determine display text and color
			var displayText string
			var textColor color.Color

			if isSelected {
				textColor = colornames.Cyan
				displayText = "[" + key.displayStr + "]"
			} else {
				if key.special != SpecialKeyNone {
					textColor = colornames.Orange
				} else {
					textColor = colornames.White
				}
				displayText = key.displayStr
			}

			// Draw the key
			// Center the text for this key
			if key.special == SpecialKeyNone {
				txtWidth := pxterm24.Font.MeasureString(displayText)
				pxterm24.Font.DrawString(
					screen,
					x-txtWidth/2,
					y,
					displayText,
					textColor,
				)
			} else {
				txtWidth := pxterm16.Font.MeasureString(displayText)
				pxterm16.Font.DrawString(
					screen,
					x-txtWidth/2,
					y+common.Pxterm16Height/5,
					displayText,
					textColor,
				)
			}
		}
	}
}
