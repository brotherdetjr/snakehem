package perftracker

import (
	"snakehem/assets/adhoc8"
	"snakehem/game/common"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

const (
	textX   = 15
	padding = 4
)

func (p *PerfTracker) Draw(screen *ebiten.Image) {
	stats := p.GetStats()
	if stats.SampleCount < 10 {
		// Wait for enough samples before displaying
		return
	}
	lines := stats.AsString()
	var lineSpacing = common.Adhoc8Height
	y := screen.Bounds().Dy() - len(lines)*lineSpacing - 2*padding
	// First line (TPS/FPS) always in light grey
	adhoc8.Font.DrawString(screen, textX, y, lines[0], colornames.Lightgrey)
	y += lineSpacing
	// Update line - red if warning, otherwise light grey
	updateColor := colornames.Lightgrey
	if stats.UpdateWarning {
		updateColor = colornames.Red
	}
	adhoc8.Font.DrawString(screen, textX, y, lines[1], updateColor)
	y += lineSpacing
	// Draw line - red if warning, otherwise light grey
	drawColor := colornames.Lightgrey
	if stats.DrawWarning {
		drawColor = colornames.Red
	}
	adhoc8.Font.DrawString(screen, textX, y, lines[2], drawColor)
}
