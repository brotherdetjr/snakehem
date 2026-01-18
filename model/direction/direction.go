package direction

type Direction uint8

const (
	None Direction = iota
	Up
	Down
	Left
	Right
)

func (d Direction) Dx() int {
	if d == Left {
		return -1
	} else if d == Right {
		return 1
	} else {
		return 0
	}
}

func (d Direction) Dy() int {
	if d == Up {
		return -1
	} else if d == Down {
		return 1
	} else {
		return 0
	}
}
