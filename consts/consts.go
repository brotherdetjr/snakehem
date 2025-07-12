package consts

import (
	"golang.org/x/image/colornames"
	"image/color"
)

const (
	Fps                           = 10
	TpsMultiplier                 = 6
	Tps                           = Fps * TpsMultiplier
	GridSize                      = 63
	CellDimPx                     = 11
	GridDimPx                     = GridSize * CellDimPx
	MaxScoresAtTop                = 5
	CountdownSeconds              = 4
	SnakeTargetLength             = 50
	HealthReductionPerBite        = 10
	NippedTailLinkBonusMultiplier = 2
	BitLinkScore                  = 1
	AppleScore                    = 50
	TargetScore                   = 999
	MaxSnakes                     = len(SnakeColours)
	ApproachingTargetScoreGap     = SnakeTargetLength * NippedTailLinkBonusMultiplier * 1.2
	GridFadeCountdown             = TpsMultiplier * 15
	NewAppleProbabilityParam      = Tps * 3
	EyeRadiusPx                   = 2
	EyeGapPx                      = 3
)

var SnakeColours = [...]color.Color{
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
