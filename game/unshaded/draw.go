package unshaded

import "github.com/hajimehoshi/ebiten/v2"

func (s *State) Draw(screen *ebiten.Image) {
	if s.perfTracker != nil {
		s.perfTracker.Draw(screen)
	}
}
