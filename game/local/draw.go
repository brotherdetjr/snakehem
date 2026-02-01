package local

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (s *State) Draw(screen *ebiten.Image) {
	if s.textInput != nil {
		s.textInput.Draw(screen)
	}
}
