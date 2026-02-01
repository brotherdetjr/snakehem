package textinput

func (t *TextInput) Submit() {
	if err := t.validation(t.value); err != nil {
		t.error = err
	} else {
		t.callback(t.value)
	}
}

func (t *TextInput) AddChar() {
	if len(t.value) < t.maxLength {
		t.value += string(AvailableChars[t.charIndex])
	}
}

func (t *TextInput) DelChar() {
	if t.value != "" {
		t.value = t.value[:len(t.value)-1]
	}
}

func (t *TextInput) NextChar() {
	t.charIndex++
	if t.charIndex == len(AvailableChars) {
		t.charIndex = 0
	}
}

func (t *TextInput) PrevChar() {
	t.charIndex--
	if t.charIndex == -1 {
		t.charIndex = len(AvailableChars) - 1
	}
}

func (t *TextInput) Update() {
	c := t.controller
	if c.IsAnyJustPressed() {
		t.error = nil
	}
	if c.IsUpJustPressed() {
		t.PrevChar()
	} else if c.IsDownJustPressed() {
		t.NextChar()
	} else if c.IsRightJustPressed() {
		t.AddChar()
	} else if c.IsLeftJustPressed() {
		t.DelChar()
	} else if c.IsStartJustPressed() {
		t.Submit()
	}
}
