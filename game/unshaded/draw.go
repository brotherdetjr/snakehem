package unshaded

import "github.com/hajimehoshi/ebiten/v2"

func (c *Content) Draw(screen *ebiten.Image) {
	if c.perfTracker != nil {
		c.perfTracker.Draw(screen)
	}
}
