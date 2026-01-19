package perception

import "snakehem/model"

type Perception struct {
	Grid      [model.GridSize][model.GridSize]any
	Countdown int
}

func NewPerception() Perception {
	return Perception{
		Grid:      [model.GridSize][model.GridSize]any{},
		Countdown: 3,
	}
}
