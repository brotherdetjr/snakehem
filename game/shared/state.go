package shared

import (
	"image/color"
	"math"
	"math/rand/v2"
	"snakehem/game/common"
	"snakehem/game/shared/scoreboard"
	"snakehem/model"
	"snakehem/model/snake"
	"snakehem/util"

	"github.com/rs/zerolog/log"
)

type State struct {
	Stage            Stage
	Grid             [model.GridSize][model.GridSize]any
	Snakes           []*snake.Snake
	Countdown        int
	FadeCountdown    int
	ActionFrameCount uint64
	scoreboard       *scoreboard.Scoreboard
	applePos         *util.Coords
}

func NewSharedState() *State {
	return &State{
		Stage:            Lobby,
		Grid:             [model.GridSize][model.GridSize]any{},
		Countdown:        3,
		FadeCountdown:    0,
		ActionFrameCount: 0,
		scoreboard:       nil,
		applePos:         nil,
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

func (s *State) TryToPutNewApple() {
	if s.applePos == nil && rand.IntN(model.NewAppleProbabilityParam) == 0 {
		x, y := s.randomUnoccupiedCell()
		if x != -1 && y != -1 {
			s.applePos = &util.Coords{X: x, Y: y}
			log.Debug().Int("x", x).Int("y", y).Msg("Put a new apple")
		}
	}
}

func (s *State) IncScore(snake *snake.Snake, delta int) {
	snake.Score += delta
	log.Debug().Int("snakeId", snake.Id).Int("score", snake.Score).Msg("New score")
	if snake.Score >= model.TargetScore {
		log.Info().Msg("Stopping the action!")
		s.FadeCountdown = model.GridFadeCountdown
	}
}

func (s *State) IsAppleHere(x, y int) bool {
	return s.applePos != nil && *s.applePos == util.Coords{X: x, Y: y}
}

func (s *State) EatApple(snake *snake.Snake) {
	s.applePos = nil
	s.IncScore(snake, model.AppleScore)
	log.Debug().Int("snakeId", snake.Id).Msg("Apple eaten!")
}

func (s *State) randomUnoccupiedCell() (int, int) {
	x := rand.IntN(model.GridSize)
	y := rand.IntN(model.GridSize)
	for ; y < model.GridSize; y++ {
		for ; x < model.GridSize; x++ {
			if s.Grid[y][x] == nil {
				return x, y
			}
		}
		x = 0
	}
	return -1, -1
}
