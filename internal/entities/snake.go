package entities

import (
	"image/color"
	"math"

	"snakehem/direction"
	"snakehem/internal/config"
)

// Snake represents a player's snake
type Snake struct {
	ID                string
	Links             []*Link
	Colour            color.Color
	Direction         direction.Direction
	Score             int
	ControllerID      string // Changed from Controller to just ID
	HeadRednessGrowth float32
}

// Link represents a single segment of a snake
type Link struct {
	HealthPercent int8
	Snake         *Snake
	X             int
	Y             int
	Redness       float32
}

// NewSnake creates a new snake with the specified controller ID and color
func NewSnake(controllerID string, colour color.Color) *Snake {
	snake := &Snake{
		ID:                controllerID,
		Links:             make([]*Link, 1),
		Colour:            colour,
		Score:             0,
		ControllerID:      controllerID,
		HeadRednessGrowth: -1,
	}
	snake.Links[0] = &Link{
		HealthPercent: 100,
		Snake:         snake,
		Redness:       1,
	}
	return snake
}

// PickInitialDirection selects an initial direction for the snake
// based on its position relative to the grid center
func (s *Snake) PickInitialDirection(gridSize int) {
	head := s.Links[0]
	x := head.X
	y := head.Y
	midPoint := gridSize/2 + 1
	dir := direction.None

	if math.Abs(float64(midPoint-x)) > math.Abs(float64(midPoint-y)) {
		if midPoint < x {
			dir = direction.Left
		} else {
			dir = direction.Right
		}
	} else {
		if midPoint < y {
			dir = direction.Up
		} else {
			dir = direction.Down
		}
	}

	s.Direction = dir
}

// ChangeRedness changes the link's redness value
func (l *Link) ChangeRedness(delta float32, tpsMultiplier int) {
	l.Redness += delta / float32(tpsMultiplier)
	if l.Redness < 0 {
		l.Redness = 0
	} else if l.Redness > 1 {
		l.Redness = 1
	}
}

// ChangeRednessWithConfig is a convenience method that uses config for TPS multiplier
func (l *Link) ChangeRednessWithConfig(delta float32, cfg *config.GameConfig) {
	l.ChangeRedness(delta, cfg.TpsMultiplier)
}
