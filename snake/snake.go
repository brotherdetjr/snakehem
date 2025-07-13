package snake

import (
	"image/color"
	"math"
	"snakehem/consts"
	"snakehem/controllers/controller"
	. "snakehem/direction"
)

type Snake struct {
	Links             []*Link
	Colour            color.Color
	Direction         Direction
	Score             int
	Controller        controller.Controller
	HeadRednessGrowth float32
}

type Link struct {
	HealthPercent int8
	Snake         *Snake
	X             int
	Y             int
	Redness       float32
}

func NewSnake(controller controller.Controller, colour color.Color) *Snake {
	snake := &Snake{
		Links:             make([]*Link, 1),
		Colour:            colour,
		Score:             0,
		Controller:        controller,
		HeadRednessGrowth: -1,
	}
	snake.Links[0] = &Link{
		HealthPercent: 100,
		Snake:         snake,
		Redness:       1,
	}
	return snake
}

func (s *Snake) PickInitialDirection() {
	head := s.Links[0]
	x := head.X
	y := head.Y
	midPoint := consts.GridSize/2 + 1
	direction := None
	if math.Abs(float64(midPoint-x)) > math.Abs(float64(midPoint-y)) {
		if midPoint < x {
			direction = Left
		} else {
			direction = Right
		}
	} else {
		if midPoint < y {
			direction = Up
		} else {
			direction = Down
		}
	}
	s.Direction = direction
}

func (l *Link) ChangeRedness(delta float32) {
	l.Redness += delta / consts.TpsMultiplier
	if l.Redness < 0 {
		l.Redness = 0
	} else if l.Redness > 1 {
		l.Redness = 1
	}
}
