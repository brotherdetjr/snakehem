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
	s.textInput = textinput.NewTextInput(
		playerName,
		"ENTER YOUR NAME",
		model.MaxNameLength,
		c,
		func(name string) {
			cb(name)
			s.stage = Off
			s.textInput = nil
		},
	).ValidateNotEmpty("name cannot be empty")
}
