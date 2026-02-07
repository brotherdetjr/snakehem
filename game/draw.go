package game

import (
	"snakehem/assets/adhoc8"
	"snakehem/game/common"
	"snakehem/model"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func (g *Game) Draw(screen *ebiten.Image) {
	if g.perfTracker != nil {
		start := time.Now()
		defer func() {
			g.perfTracker.RecordDraw(time.Since(start))
			g.perfTracker.RecordFPS(ebiten.ActualFPS())
		}()
	}

	if ebiten.Tick()%model.TpsMultiplier == 0 {
		g.lastFrame.Clear()
		g.sharedState.Draw(g.lastFrame)
		g.localState.Draw(g.lastFrame)
		g.applyShader(g.lastFrame)
		if g.perfTrackerVisible {
			g.drawPerfStats(g.lastFrame)
		}
	}
	screen.DrawImage(g.lastFrame, nil)
}

func (g *Game) drawPerfStats(screen *ebiten.Image) {
	stats := g.perfTracker.GetStats()
	if stats.SampleCount < 10 {
		// Wait for enough samples before displaying
		return
	}

	lines := stats.AsString()
	y := 20
	for _, line := range lines {
		adhoc8.Font.DrawString(screen, 15, y, line, colornames.Lightgray)
		y += common.Adhoc8Height + 1
	}
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
