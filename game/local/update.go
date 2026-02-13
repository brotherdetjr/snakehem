package local

import (
	"snakehem/game/common"
)

func (s *State) Update(ctx *common.Context) {
	if s.textInput != nil {
		s.textInput.Update(ctx)
	}
}
