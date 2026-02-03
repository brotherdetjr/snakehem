package local

import (
	"image/color"
	"snakehem/game/local/textinput"
	"snakehem/input/controller"
	"snakehem/model"
)

func (s *State) Update() {
	if s.textInput != nil {
		s.textInput.Update()
	}
}

func (s *State) StagePlayerName(c controller.Controller, playerName string, colour color.Color, cb func(string)) {
	s.stage = PlayerName
	s.textInput = textinput.
		NewTextInput(c).
		WithLabel("ENTER YOUR NAME").
		WithValue(playerName).
		WithMaxLength(model.MaxNameLength).
		WithTextColour(colour).
		WithKeyboardCols(12).
		WithNameCapitalisation(true).
		ValidateNotEmpty("name cannot be empty").
		WithCallback(func(name string) {
			cb(name)
			s.stage = Off
			s.textInput = nil
		})
}
