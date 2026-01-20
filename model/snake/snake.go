package snake

import (
	"image/color"
	"math"
	"snakehem/model"
	"snakehem/model/direction"
)

type Snake struct {
	Id                int
	Links             []*Link
	Colour            color.Color
	Direction         direction.Direction
	Score             int
	HeadRednessGrowth float32
}

type Link struct {
	SnakeId       int
	HealthPercent int8
	X             int
	Y             int
	Redness       float32
}

func NewSnake(id int, colour color.Color) *Snake {
	snake := &Snake{
		Id:                id,
		Links:             make([]*Link, 1),
		Colour:            colour,
		Score:             0,
		HeadRednessGrowth: -1,
	}
	snake.Links[0] = &Link{
		HealthPercent: 100,
		SnakeId:       id,
		Redness:       1,
	}
	return snake
}

func (s *Snake) PickInitialDirection() {
	head := s.Links[0]
	x := head.X
	y := head.Y
	midPoint := model.GridSize/2 + 1
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

func (l *Link) ChangeRedness(delta float32) {
	l.Redness += delta / model.TpsMultiplier
	if l.Redness < 0 {
		l.Redness = 0
	} else if l.Redness > 1 {
		l.Redness = 1
	}
}
