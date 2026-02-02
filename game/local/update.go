package local

import (
	"snakehem/game/local/textinput"
	"snakehem/input/controller"
	"snakehem/model"
)

func (s *State) Update() {
	if s.textInput != nil {
		s.textInput.Update()
	}
}

func (s *State) StagePlayerName(c controller.Controller, playerName string, cb func(string)) {
	s.stage = PlayerName
	s.textInput = textinput.
		NewTextInput(c).
		WithLabel("ENTER YOUR NAME").
		WithValue(playerName).
		WithMaxLength(model.MaxNameLength).
		ValidateNotEmpty("name cannot be empty").
		WithCallback(func(name string) {
			cb(name)
			s.stage = Off
			s.textInput = nil
		})
}
