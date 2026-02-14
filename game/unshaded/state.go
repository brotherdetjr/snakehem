package unshaded

import (
	"snakehem/game/unshaded/perftracker"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Content struct {
	perfTracker *perftracker.PerfTracker
}

func NewContent() *Content {
	return &Content{perfTracker: nil}
}

type Stage uint8

func (c *Content) RecordUpdateTimeAndTps(since time.Time) {
	if c.perfTracker != nil {
		c.perfTracker.RecordUpdate(time.Since(since))
		c.perfTracker.RecordTPS(ebiten.ActualTPS())
	}
}

func (c *Content) RecordDrawTimeAndFps(since time.Time) {
	if c.perfTracker != nil {
		c.perfTracker.RecordDraw(time.Since(since))
		c.perfTracker.RecordFPS(ebiten.ActualFPS())
	}
}
