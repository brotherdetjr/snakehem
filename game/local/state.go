package local

import "snakehem/game/local/textinput"

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
