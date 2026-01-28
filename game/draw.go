package game

import (
	"snakehem/game/common"
	"snakehem/game/shared"
	"snakehem/model"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Draw(screen *ebiten.Image) {
	if g.frameCount%model.TpsMultiplier == 0 {
		// Render shared state
		// TODO: check if different from lastState, after server-client is implemented
		g.sharedFrame.Clear()
		shared.DrawSharedState(g.sharedState, g.sharedFrame)

		if g.localState.Dirty() {
			g.localFrame.Clear()
			g.localState.DrawLocalState(g.localFrame)
		}
	}

	// Composite shared and local frames
	composite := ebiten.NewImage(common.GridDimPx, common.GridDimPx)
	composite.DrawImage(g.sharedFrame, nil)
	composite.DrawImage(g.localFrame, nil)

	// Apply shader to composite
	g.applyShader(composite)
	screen.DrawImage(composite, nil)
}

func (g *Game) applyShader(screen *ebiten.Image) {
	w := screen.Bounds().Dx()
	h := screen.Bounds().Dy()
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = screen
	opts.Uniforms = map[string]interface{}{
		// Kage uniforms here
	}
	img := ebiten.NewImage(w, h)
	img.DrawRectShader(w, h, g.shader, opts)
	screen.DrawImage(img, nil)
}
