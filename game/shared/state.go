package shared

import (
	"image/color"
	"math"
	"math/rand/v2"
	"snakehem/game/common"
	"snakehem/game/shared/scoreboard"
	"snakehem/game/shared/snake"
	"snakehem/model"
	"snakehem/util"

	"github.com/rs/zerolog/log"
)

type Content struct {
	Stage            Stage
	Grid             [model.GridSize][model.GridSize]any
	Snakes           []*snake.Snake
	Countdown        int
	FadeCountdown    int
	ActionFrameCount uint64
	scoreboard       *scoreboard.Scoreboard
	applePos         *util.Coords
}

func NewContent() *Content {
	return &Content{
		Stage:            Lobby,
		Grid:             [model.GridSize][model.GridSize]any{},
		Countdown:        model.Tps * model.CountdownSeconds,
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

func (c *Content) SwitchToScoreboardStage() {
	c.Stage = Scoreboard
	entries := make([]scoreboard.Entry, len(c.Snakes))
	for i, s := range c.Snakes {
		score := s.Score
		if score > model.TargetScore {
			score = model.TargetScore
		}
		entries[i] = scoreboard.Entry{
			Name:  s.Name,
			Score: score,
			ColourFunc: func() color.Color {
				return common.WithRedness(s.Colour, s.Links[0].Redness)
			},
		}
	}
	c.scoreboard = scoreboard.NewScoreboard(entries)
}

func (c *Content) SwitchToLobbyStage() {
	c.Stage = Lobby
	c.Grid = [model.GridSize][model.GridSize]any{}
	for _, s := range c.Snakes {
		s.Score = 0
		s.Links = s.Links[0:1]
	}
	c.Countdown = model.Tps * model.CountdownSeconds
	c.FadeCountdown = 0
	c.ActionFrameCount = 0
	c.scoreboard = nil
	c.applePos = nil
	c.LayoutSnakes()
	log.Info().Msg("Game restarted")
}

func (c *Content) LayoutSnakes() {
	delta := 2 * math.Pi / float64(len(c.Snakes))
	alpha := float64(0)
	for _, s := range c.Snakes {
		y := model.GridSize/2 - int(math.Cos(alpha)*model.GridSize/3)
		x := model.GridSize/2 + int(math.Sin(alpha)*model.GridSize/3)
		head := s.Links[0]
		head.X = x
		head.Y = y
		c.Grid[y][x] = head
		alpha += delta
		s.PickInitialDirection()
	}
}

func (c *Content) TryToPutNewApple() {
	if c.applePos == nil && rand.IntN(model.NewAppleProbabilityParam) == 0 {
		x, y := c.randomUnoccupiedCell()
		if x != -1 && y != -1 {
			c.applePos = &util.Coords{X: x, Y: y}
			log.Debug().Int("x", x).Int("y", y).Msg("Put a new apple")
		}
	}
}

func (c *Content) IncScore(snake *snake.Snake, delta int) {
	snake.Score += delta
	log.Debug().Int("snakeId", snake.Id).Int("score", snake.Score).Msg("New score")
	if snake.Score >= model.TargetScore {
		log.Info().Msg("Stopping the action!")
		c.FadeCountdown = model.GridFadeCountdown
	}
}

func (c *Content) IsAppleHere(x, y int) bool {
	return c.applePos != nil && *c.applePos == util.Coords{X: x, Y: y}
}

func (c *Content) EatApple(snake *snake.Snake) {
	c.applePos = nil
	c.IncScore(snake, model.AppleScore)
	log.Debug().Int("snakeId", snake.Id).Msg("Apple eaten!")
}

func (c *Content) GetCountdownSeconds() int {
	return (c.Countdown - 1) / model.Tps
}

func (c *Content) randomUnoccupiedCell() (int, int) {
	x := rand.IntN(model.GridSize)
	y := rand.IntN(model.GridSize)
	for ; y < model.GridSize; y++ {
		for ; x < model.GridSize; x++ {
			if c.Grid[y][x] == nil {
				return x, y
			}
		}
		x = 0
	}
	return -1, -1
}
