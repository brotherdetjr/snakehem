package snake

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	. "snakehem/direction"
)

type Snake struct {
	Links             []*Link
	Colour            color.Color
	Direction         Direction
	Score             int
	Id                ebiten.GamepadID
	HeadRednessGrowth float32
}

type Link struct {
	HealthPercent int8
	Snake         *Snake
	X             int
	Y             int
	Redness       float32
}

func NewSnake(id ebiten.GamepadID, colour color.Color) *Snake {
	snake := &Snake{
		Links:             make([]*Link, 1),
		Colour:            colour,
		Direction:         None,
		Score:             0,
		Id:                id,
		HeadRednessGrowth: -1,
	}
	snake.Links[0] = &Link{
		HealthPercent: 100,
		Snake:         snake,
		Redness:       1,
	}
	return snake
}

func (l *Link) ChangeRedness(delta float32) {
	l.Redness += delta
	if l.Redness < 0 {
		l.Redness = 0
	} else if l.Redness > 1 {
		l.Redness = 1
	}
}
