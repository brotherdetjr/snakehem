package textinput

func (t *TextInput) Submit() {
	if err := t.validation(t.value); err != nil {
		t.error = err
	} else {
		t.callback(t.value)
	}
}

func (t *TextInput) AddChar() {
	// Legacy method - now uses cursor position
	t.AddCharAtCursor()
}

func (t *TextInput) DelChar() {
	if t.value != "" {
		t.value = t.value[:len(t.value)-1]
	}
}

func (t *TextInput) moveLeft() {
	newCol := t.cursorCol - 1
	newRow := t.cursorRow

	// Wrap to end of current row
	if newCol < 0 {
		newCol = KeyboardCols - 1
	}

	// If position invalid, keep searching left with wrapping
	startCol := newCol
	for !isValidPosition(newRow, newCol) {
		newCol--
		if newCol < 0 {
			newCol = KeyboardCols - 1
		}
		if newCol == startCol {
			break // Prevent infinite loop
		}
	}

	if isValidPosition(newRow, newCol) {
		t.cursorRow = newRow
		t.cursorCol = newCol
	}
}

func (t *TextInput) moveRight() {
	newCol := t.cursorCol + 1
	newRow := t.cursorRow

	// Wrap to start of current row
	if newCol >= KeyboardCols {
		newCol = 0
	}

	// If position invalid, keep searching right with wrapping
	startCol := newCol
	for !isValidPosition(newRow, newCol) {
		newCol++
		if newCol >= KeyboardCols {
			newCol = 0
		}
		if newCol == startCol {
			break // Prevent infinite loop
		}
	}

	if isValidPosition(newRow, newCol) {
		t.cursorRow = newRow
		t.cursorCol = newCol
	}
}

func (t *TextInput) moveUp() {
	newRow := t.cursorRow - 1
	newCol := t.cursorCol

	// Wrap to bottom
	if newRow < 0 {
		newRow = KeyboardRows - 1
	}

	// If position invalid, try to find valid position
	startRow := newRow
	for !isValidPosition(newRow, newCol) {
		newRow--
		if newRow < 0 {
			newRow = KeyboardRows - 1
		}
		if newRow == startRow {
			break // Prevent infinite loop
		}
	}

	if isValidPosition(newRow, newCol) {
		t.cursorRow = newRow
		t.cursorCol = newCol
	}
}

func (t *TextInput) moveDown() {
	newRow := t.cursorRow + 1
	newCol := t.cursorCol

	// Wrap to top
	if newRow >= KeyboardRows {
		newRow = 0
	}

	// If position invalid, try to find valid position
	startRow := newRow
	for !isValidPosition(newRow, newCol) {
		newRow++
		if newRow >= KeyboardRows {
			newRow = 0
		}
		if newRow == startRow {
			break // Prevent infinite loop
		}
	}

	if isValidPosition(newRow, newCol) {
		t.cursorRow = newRow
		t.cursorCol = newCol
	}
}

func (t *TextInput) AddCharAtCursor() {
	key := t.getCurrentKey()
	if key != nil && key.char != 0 && len(t.value) < t.maxLength {
		t.value += string(key.char)
	}
}

func (t *TextInput) pressCurrentKey() {
	key := t.getCurrentKey()
	if key == nil {
		return
	}

	switch key.special {
	case SpecialKeyOK:
		t.Submit()
	case SpecialKeyDEL:
		t.DelChar()
	default:
		// Regular character
		t.AddCharAtCursor()
	}
}

func (t *TextInput) Update() {
	c := t.controller
	if c.IsAnyJustPressed() {
		t.error = nil
	}
	if c.IsUpJustPressed() {
		t.moveUp()
	} else if c.IsDownJustPressed() {
		t.moveDown()
	} else if c.IsLeftJustPressed() {
		t.moveLeft()
	} else if c.IsRightJustPressed() {
		t.moveRight()
	} else if c.IsStartJustPressed() {
		t.pressCurrentKey()
	}
}
