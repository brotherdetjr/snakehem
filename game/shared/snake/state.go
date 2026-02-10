package snake

import (
	"image/color"
	"math"
	"snakehem/model"
)

type Snake struct {
	Id        int
	Name      string
	Links     []*Link
	Colour    color.Color
	Direction Direction
	Score     int
}

type Link struct {
	SnakeId       int
	HealthPercent int8
	X             int
	Y             int
	Redness       float32
}

func NewSnake(id int, name string, colour color.Color) *Snake {
	snake := &Snake{
		Id:     id,
		Links:  make([]*Link, 1),
		Colour: colour,
		Score:  0,
		Name:   name,
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
	dir := None
	if math.Abs(float64(midPoint-x)) > math.Abs(float64(midPoint-y)) {
		if midPoint < x {
			dir = Left
		} else {
			dir = Right
		}
	} else {
		if midPoint < y {
			dir = Up
		} else {
			dir = Down
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

type Direction uint8

const (
	None Direction = iota
	Up
	Down
	Left
	Right
)

func (d Direction) Dx() int {
	if d == Left {
		return -1
	} else if d == Right {
		return 1
	} else {
		return 0
	}
}

func (d Direction) Dy() int {
	if d == Up {
		return -1
	} else if d == Down {
		return 1
	} else {
		return 0
	}
}
