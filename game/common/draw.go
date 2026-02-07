package common

import (
	"fmt"
	"image/color"
	"math"
	"snakehem/assets/adhoc8"
	"snakehem/assets/pxterm24"
	"snakehem/model"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pbnjay/pixfont"
	"golang.org/x/image/colornames"
)

var ScoreFmt = "%0" + fmt.Sprint(int(math.Log10(model.TargetScore))+1) + "d"
var Pxterm16Height = pxterm24.Font.GetHeight()
var Pxterm24Height = pxterm24.Font.GetHeight()
var Adhoc8Height = adhoc8.Font.GetHeight()

const (
	CellDimPx = 11
	GridDimPx = model.GridSize * CellDimPx
)

var SnakeColours = [model.MaxSnakes]color.Color{
	colornames.Lightgrey,
	color.NRGBA{ // orange
		R: 255,
		G: 128,
		B: 10,
		A: 255,
	},
	colornames.Yellow,
	color.NRGBA{ // dark green
		R: 100,
		G: 170,
		B: 0,
		A: 255,
	},
	colornames.Cyan,
	colornames.Blue,
	color.NRGBA{ // dark blue
		R: 0,
		G: 0,
		B: 100,
		A: 255,
	},
	color.NRGBA{ // dark magenta
		R: 100,
		G: 0,
		B: 84,
		A: 255,
	},
	colornames.Magenta,
}

func DrawTextCentered(screen *ebiten.Image, txt string, colour color.Color, top float64, font *pixfont.PixFont) {
	txtWidth := font.MeasureString(txt)
	font.DrawString(screen, (GridDimPx-txtWidth)/2, int(top), txt, colour)
}

// WithRedness transforms a given colour by add a red hue to it. redness argument
// varies from 0 (keep the original colour) to 1 (make it fully red).
func WithRedness(colour color.Color, redness float32) color.Color {
	if redness < 0 || redness > 1 {
		panic("redness must be between 0 and 1")
	}
	red, green, blue, _ := colour.RGBA()
	r := float32(red >> 8)
	g := float32(green >> 8)
	b := float32(blue >> 8)
	return color.NRGBA{
		R: uint8(r + (255-r)*redness),
		G: uint8(g * (1 - redness)),
		B: uint8(b * (1 - redness)),
		A: 255,
	}
}
