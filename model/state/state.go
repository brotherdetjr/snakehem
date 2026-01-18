package state

type State uint8

const (
	Lobby State = iota
	Action
	Scoreboard
)
