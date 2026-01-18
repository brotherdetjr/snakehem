package graphics

import (
	"image/color"
	"snakehem/model"

	"golang.org/x/image/colornames"
)

const (
	CellDimPx      = 11
	GridDimPx      = model.GridSize * CellDimPx
	MaxScoresAtTop = 5
	EyeRadiusPx    = 2
	EyeGapPx       = 3
)

var SnakeColours = [model.MaxSnakes]color.Color{
	colornames.Lightgrey,
	color.NRGBA{
		R: 255,
		G: 128,
		B: 10,
		A: 255,
	},
	colornames.Yellow,
	color.NRGBA{
		R: 100,
		G: 170,
		B: 0,
		A: 255,
	},
	colornames.Cyan,
	colornames.Blue,
	color.NRGBA{
		R: 0,
		G: 0,
		B: 100,
		A: 255,
	},
	color.NRGBA{
		R: 100,
		G: 0,
		B: 84,
		A: 255,
	},
	colornames.Magenta,
}
