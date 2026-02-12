package game

import (
	"snakehem/model"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Draw(screen *ebiten.Image) {
	start := time.Now()
	defer func() {
		g.localState.RecordDrawTimeAndFps(time.Since(start), ebiten.ActualFPS())
	}()

	if ebiten.Tick()%model.TpsMultiplier == 0 {
		g.lastFrame.Clear()
		g.sharedState.Draw(g.lastFrame)
		g.localState.Draw(g.lastFrame)
		g.applyShader(g.lastFrame)
	}
	screen.DrawImage(g.lastFrame, nil)
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
