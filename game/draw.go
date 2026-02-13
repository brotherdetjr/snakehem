package game

import (
	"snakehem/game/common"
	"snakehem/model"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Draw(screen *ebiten.Image) {
	if ebiten.Tick()%model.TpsMultiplier == 0 {
		g.doDraw(screen)
	}
}

func (g *Game) doDraw(screen *ebiten.Image) {
	start := time.Now()
	defer func() {
		g.unshadedState.RecordDrawTimeAndFps(start)
	}()
	frame := ebiten.NewImage(common.GridDimPx, common.GridDimPx)
	g.sharedState.Draw(frame)
	g.localState.Draw(frame)
	g.applyShader(frame)
	g.unshadedState.Draw(frame)
	screen.DrawImage(frame, nil)
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
