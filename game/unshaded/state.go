package unshaded

import (
	"snakehem/game/unshaded/perftracker"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type State struct {
	perfTracker *perftracker.PerfTracker
}

func NewUnshadedState() *State {
	return &State{perfTracker: nil}
}

type Stage uint8

func (s *State) RecordUpdateTimeAndTps(since time.Time) {
	if s.perfTracker != nil {
		s.perfTracker.RecordUpdate(time.Since(since))
		s.perfTracker.RecordTPS(ebiten.ActualTPS())
	}
}

func (s *State) RecordDrawTimeAndFps(since time.Time) {
	if s.perfTracker != nil {
		s.perfTracker.RecordDraw(time.Since(since))
		s.perfTracker.RecordFPS(ebiten.ActualFPS())
	}
}
