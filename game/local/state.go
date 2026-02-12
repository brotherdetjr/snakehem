package local

import (
	"image/color"
	"snakehem/game/local/perftracker"
	"snakehem/game/local/textinput"
	"snakehem/input/controller"
	"snakehem/model"
	"time"
)

type State struct {
	stage       Stage
	textInput   *textinput.TextInput
	perfTracker *perftracker.PerfTracker
}

func NewLocalState() *State {
	return &State{
		stage:       Off,
		textInput:   nil,
		perfTracker: nil,
	}
}

func (s *State) GetStage() Stage {
	return s.stage
}

type Stage uint8

const (
	Off Stage = iota
	PlayerName
)

func (s *State) SwitchToPlayerNameStage(c controller.Controller, playerName string, colour color.Color, cb func(string)) {
	s.stage = PlayerName
	s.textInput = textinput.
		NewTextInput(c).
		WithLabel("ENTER YOUR NAME").
		WithValue(playerName).
		WithMaxLength(model.MaxNameLength).
		WithTextColour(colour).
		WithKeyboardCols(12).
		WithCapsBehaviour(textinput.CapsBehaviourNames).
		ValidateNotEmpty("name cannot be empty").
		WithCallback(func(name string) {
			cb(name)
			s.stage = Off
			s.textInput = nil
		})
}

func (s *State) RecordUpdateTimeAndTps(duration time.Duration, tps float64) {
	if s.perfTracker != nil {
		s.perfTracker.RecordUpdate(duration)
		s.perfTracker.RecordTPS(tps)
	}
}

func (s *State) RecordDrawTimeAndFps(duration time.Duration, fps float64) {
	if s.perfTracker != nil {
		s.perfTracker.RecordDraw(duration)
		s.perfTracker.RecordFPS(fps)
	}
}
