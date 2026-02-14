package unshaded

import (
	"snakehem/game/unshaded/perftracker"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/rs/zerolog/log"
)

func (c *Content) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		if c.perfTracker == nil {
			c.perfTracker = perftracker.NewPerfTracker()
		} else {
			c.perfTracker = nil
		}
		log.Info().Bool("enabled", c.perfTracker != nil).Msg("Performance tracker")
	}
}
