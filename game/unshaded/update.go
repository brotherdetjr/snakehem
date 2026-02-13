package unshaded

import (
	"snakehem/game/unshaded/perftracker"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/rs/zerolog/log"
)

func (s *State) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		if s.perfTracker == nil {
			s.perfTracker = perftracker.NewPerfTracker()
		} else {
			s.perfTracker = nil
		}
		log.Info().Bool("enabled", s.perfTracker != nil).Msg("Performance tracker")
	}
}
