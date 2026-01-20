package perception

import (
	"snakehem/model"
	"snakehem/model/snake"
	"snakehem/model/stage"
)

type Perception struct {
	Stage         stage.Stage
	Grid          [model.GridSize][model.GridSize]any
	Snakes        []*snake.Snake
	Countdown     int
	FadeCountdown int
}

func NewPerception() Perception {
	return Perception{
		Stage:         stage.Lobby,
		Grid:          [model.GridSize][model.GridSize]any{},
		Countdown:     3,
		FadeCountdown: 0,
	}
}
