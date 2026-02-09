package shared

import (
	"image/color"
	"math"
	"snakehem/game/common"
	"snakehem/game/shared/scoreboard"
	"snakehem/model"
	"snakehem/model/snake"
)

type State struct {
	Stage            Stage
	Grid             [model.GridSize][model.GridSize]any
	Snakes           []*snake.Snake
	Countdown        int
	FadeCountdown    int
	ActionFrameCount uint64
	scoreboard       *scoreboard.Scoreboard
}

func NewSharedState() *State {
	return &State{
		Stage:            Lobby,
		Grid:             [model.GridSize][model.GridSize]any{},
		Countdown:        3,
		FadeCountdown:    0,
		ActionFrameCount: 0,
		scoreboard:       nil,
	}
}

type Stage uint8

const (
	Lobby Stage = iota
	Action
	Scoreboard
)

func (s *State) SwitchToScoreboardStage() {
	s.Stage = Scoreboard
	entries := make([]scoreboard.Entry, len(s.Snakes))
	for i, sn := range s.Snakes {
		score := sn.Score
		if score > model.TargetScore {
			score = model.TargetScore
		}
		entries[i] = scoreboard.Entry{
			Name:  sn.Name,
			Score: score,
			ColourFunc: func() color.Color {
				return common.WithRedness(sn.Colour, sn.Links[0].Redness)
			},
		}
	}
	s.scoreboard = scoreboard.NewScoreboard(entries)
}

func (s *State) SwitchToLobbyStage() {
	s.Grid = [model.GridSize][model.GridSize]any{}
	s.Stage = Lobby
	s.FadeCountdown = 0
	s.ActionFrameCount = 0
	for _, sn := range s.Snakes {
		sn.Score = 0
		sn.Links = sn.Links[0:1]
	}
	s.LayoutSnakes()
}

func (s *State) LayoutSnakes() {
	delta := 2 * math.Pi / float64(len(s.Snakes))
	alpha := float64(0)
	for _, sn := range s.Snakes {
		y := model.GridSize/2 - int(math.Cos(alpha)*model.GridSize/3)
		x := model.GridSize/2 + int(math.Sin(alpha)*model.GridSize/3)
		head := sn.Links[0]
		head.X = x
		head.Y = y
		s.Grid[y][x] = head
		alpha += delta
		sn.PickInitialDirection()
	}
}
