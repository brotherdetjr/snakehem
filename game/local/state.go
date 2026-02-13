package local

import (
	"image/color"
	"snakehem/game/local/textinput"
	"snakehem/input/controller"
	"snakehem/model"
)

type State struct {
	stage     Stage
	textInput *textinput.TextInput
}

func NewLocalState() *State {
	return &State{
		stage:     Off,
		textInput: nil,
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
	if s.textInput != nil {
		return
	}
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
