package perftracker

import (
	"snakehem/assets/adhoc8"
	"snakehem/game/common"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func (p *PerfTracker) Draw(screen *ebiten.Image) {
	stats := p.GetStats()
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
