package local

import (
	"image/color"
	"snakehem/game/local/textinput"
	"snakehem/input/controller"
	"snakehem/model"
)

type Content struct {
	stage     Stage
	textInput *textinput.TextInput
}

func NewContent() *Content {
	return &Content{
		stage:     Off,
		textInput: nil,
	}
}

func (c *Content) GetStage() Stage {
	return c.stage
}

type Stage uint8

const (
	Off Stage = iota
	PlayerName
)

func (c *Content) SwitchToPlayerNameStage(ctrl controller.Controller, playerName string, colour color.Color, cb func(string)) {
	if c.textInput != nil {
		return
	}
	c.stage = PlayerName
	c.textInput = textinput.
		NewTextInput(ctrl).
		WithLabel("ENTER YOUR NAME").
		WithValue(playerName).
		WithMaxLength(model.MaxNameLength).
		WithTextColour(colour).
		WithKeyboardCols(12).
		WithCapsBehaviour(textinput.CapsBehaviourNames).
		ValidateNotEmpty("name cannot be empty").
		WithCallback(func(name string) {
			cb(name)
			c.stage = Off
			c.textInput = nil
		})
}
