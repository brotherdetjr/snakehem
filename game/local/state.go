package local

type State struct {
	textInput *TextInput
}

func NewLocalState() *State {
	return &State{
		textInput: nil,
	}
}

func (s *State) SetTextInput(textInput *TextInput) {
	s.textInput = textInput
}

func (s *State) Update() {
	if s.textInput != nil {
		s.textInput.Update()
	}
}
