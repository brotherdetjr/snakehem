package local

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (c *Content) Draw(screen *ebiten.Image) {
	if c.textInput != nil {
		c.textInput.Draw(screen)
	}
}
