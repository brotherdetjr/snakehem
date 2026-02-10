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

	// First line (TPS/FPS) always in light grey
	adhoc8.Font.DrawString(screen, 15, y, lines[0], colornames.Lightgrey)
	y += common.Adhoc8Height + 1

	// Update line - red if warning, otherwise light grey
	updateColor := colornames.Lightgrey
	if stats.UpdateWarning {
		updateColor = colornames.Red
	}
	adhoc8.Font.DrawString(screen, 15, y, lines[1], updateColor)
	y += common.Adhoc8Height + 1

	// Draw line - red if warning, otherwise light grey
	drawColor := colornames.Lightgrey
	if stats.DrawWarning {
		drawColor = colornames.Red
	}
	adhoc8.Font.DrawString(screen, 15, y, lines[2], drawColor)
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
