package config

import (
	"image/color"

	"golang.org/x/image/colornames"
)

// GameConfig holds all game configuration
type GameConfig struct {
	// Gameplay timing
	GameSpeedFPS  int
	TpsMultiplier int

	// Grid settings
	GridSize  int
	CellDimPx int

	// UI settings
	MaxScoresAtTop int

	// Game timing
	CountdownSeconds         int
	GridFadeCountdown        int
	NewAppleProbabilityParam int

	// Snake settings
	SnakeTargetLength int
	MaxSnakes         int

	// Scoring
	HealthReductionPerBite        int8
	NippedTailLinkBonusMultiplier int
	BitLinkScore                  int
	AppleScore                    int
	TargetScore                   int
	ApproachingTargetScoreGap     int

	// Rendering
	EyeRadiusPx float32
	EyeGapPx    float32

	// Colors
	SnakeColors []color.Color
}

// DefaultConfig returns the default game configuration
// Values match those from consts/consts.go
func DefaultConfig() *GameConfig {
	snakeColors := []color.Color{
		colornames.Lightgrey,
		color.NRGBA{R: 255, G: 128, B: 10, A: 255},
		colornames.Yellow,
		color.NRGBA{R: 100, G: 170, B: 0, A: 255},
		colornames.Cyan,
		colornames.Blue,
		color.NRGBA{R: 0, G: 0, B: 100, A: 255},
		color.NRGBA{R: 100, G: 0, B: 84, A: 255},
		colornames.Magenta,
	}

	cfg := &GameConfig{
		GameSpeedFPS:                  10,
		TpsMultiplier:                 6,
		GridSize:                      63,
		CellDimPx:                     11,
		MaxScoresAtTop:                5,
		CountdownSeconds:              4,
		SnakeTargetLength:             50,
		HealthReductionPerBite:        10,
		NippedTailLinkBonusMultiplier: 2,
		BitLinkScore:                  1,
		AppleScore:                    45,
		TargetScore:                   999,
		EyeRadiusPx:                   2,
		EyeGapPx:                      3,
		SnakeColors:                   snakeColors,
	}

	// Computed values
	cfg.MaxSnakes = len(cfg.SnakeColors)
	cfg.ApproachingTargetScoreGap = cfg.SnakeTargetLength*cfg.NippedTailLinkBonusMultiplier - 1
	cfg.GridFadeCountdown = cfg.Tps() * 15
	cfg.NewAppleProbabilityParam = cfg.Tps() * 3

	return cfg
}

// Tps returns the total ticks per second
func (c *GameConfig) Tps() int {
	return c.GameSpeedFPS * c.TpsMultiplier
}

// GridDimPx returns the grid dimension in pixels
func (c *GameConfig) GridDimPx() int {
	return c.GridSize * c.CellDimPx
}
