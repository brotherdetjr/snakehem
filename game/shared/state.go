package shared

import (
	"snakehem/model"
	"snakehem/model/snake"
	"snakehem/model/stage"
)

type State struct {
	Stage            stage.Stage
	Grid             [model.GridSize][model.GridSize]any
	Snakes           []*snake.Snake
	Countdown        int
	FadeCountdown    int
	ActionFrameCount uint64
}

func NewSharedState() *State {
	return &State{
		Stage:            stage.Lobby,
		Grid:             [model.GridSize][model.GridSize]any{},
		Countdown:        3,
		FadeCountdown:    0,
		ActionFrameCount: 0,
	}
}
