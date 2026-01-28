package local

type State struct {
	textInput *TextInput
	dirty     bool
}

func NewLocalState() *State {
	return &State{
		textInput: nil,
		dirty:     true,
	}
}

func (s *State) Dirty() bool {
	var result = s.dirty
	if s.textInput != nil && s.textInput.Dirty() {
		result = true
	}
	s.dirty = false
	return result
}

func (s *State) SetTextInput(textInput *TextInput) {
	s.textInput = textInput
	s.dirty = true
}

func (s *State) Update() {
	if s.textInput != nil {
		s.textInput.Update()
	}
}
