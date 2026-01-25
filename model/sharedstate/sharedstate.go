package sharedstate

import (
	"snakehem/model"
	"snakehem/model/snake"
	"snakehem/model/stage"
)

type SharedState struct {
	Stage         stage.Stage
	Grid          [model.GridSize][model.GridSize]any
	Snakes        []*snake.Snake
	Countdown     int
	FadeCountdown int
	ElapsedFrames uint64
}

func NewSharedState() SharedState {
	return SharedState{
		Stage:         stage.Lobby,
		Grid:          [model.GridSize][model.GridSize]any{},
		Countdown:     3,
		FadeCountdown: 0,
		ElapsedFrames: 0,
	}
}
