package game

import (
	"math"
	"math/rand/v2"
	"os"
	"slices"
	"snakehem/graphics"
	"snakehem/input"
	"snakehem/input/keyboard"
	"snakehem/model"
	. "snakehem/model/apple"
	. "snakehem/model/direction"
	. "snakehem/model/snake"
	. "snakehem/model/stage"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/rs/zerolog/log"
)

func (g *Game) Update() error {
	if keyboard.Instance.IsExitJustPressed() {
		os.Exit(0)
	}
	switch g.perception.Stage {
	case Lobby:
		g.updateHeadCount()
	case Action:
		if g.countdown > 0 {
			g.countdown--
		}
		if g.perception.FadeCountdown > 0 {
			g.perception.FadeCountdown--
			if g.perception.FadeCountdown == 0 {
				g.perception.Stage = Scoreboard
				break
			}
		}
		for _, snake := range g.perception.Snakes {
			head := snake.Links[0]
			if g.countdown <= model.Tps {
				head.ChangeRedness(0.2 * g.snakeHeadsRednessGrowth)
			} else if g.snakeControllers[snake.Id].IsAnyJustPressed() && g.perception.FadeCountdown == 0 {
				head.Redness = 1
			} else {
				head.ChangeRedness(-0.1)
			}
			for _, link := range snake.Links {
				if link != head {
					link.ChangeRedness(-0.1)
				}
			}
		}
		if g.perception.ElapsedFrames%model.Tps == 0 {
			g.snakeHeadsRednessGrowth *= -1
		}
		if g.countdown > model.Tps {
			break
		}
		for _, snake := range g.perception.Snakes {
			direction := snake.Direction
			controller := g.snakeControllers[snake.Id]
			if g.perception.FadeCountdown == 0 {
				if controller.IsUpJustPressed() {
					direction = Up
					log.Info().Any("snakeId", controller).Str("direction", "Up").Msg("New direction")
				} else if controller.IsDownJustPressed() {
					direction = Down
					log.Info().Any("snakeId", controller).Str("direction", "Down").Msg("New direction")
				} else if controller.IsLeftJustPressed() {
					direction = Left
					log.Info().Any("snakeId", controller).Str("direction", "Left").Msg("New direction")
				} else if controller.IsRightJustPressed() {
					direction = Right
					log.Info().Any("snakeId", controller).Str("direction", "Right").Msg("New direction")
				}
			}
			nX, nY := newHeadCoords(snake, direction)
			// not biting self in the neck, preserving same direction if the case
			if len(snake.Links) > 1 && nX == snake.Links[1].X && nY == snake.Links[1].Y {
				direction = snake.Direction
				nX, nY = newHeadCoords(snake, direction)
			}
			if g.perception.ElapsedFrames%model.TpsMultiplier == 0 {
				if g.perception.Grid[nY][nX] == nil {
					tail := snake.Links[len(snake.Links)-1]
					oldTailX := tail.X
					oldTailY := tail.Y
					for i := len(snake.Links) - 1; i > 0; i-- {
						link := snake.Links[i]
						prevLink := snake.Links[i-1]
						link.X = prevLink.X
						link.Y = prevLink.Y
					}
					snake.Links[0].X = nX
					snake.Links[0].Y = nY
					if len(snake.Links) < model.SnakeTargetLength {
						snake.Links = append(snake.Links, &Link{
							HealthPercent: 100,
							SnakeId:       snake.Id,
							X:             oldTailX,
							Y:             oldTailY,
							Redness:       0,
						})
					} else {
						g.perception.Grid[oldTailY][oldTailX] = nil
					}
					for _, link := range snake.Links {
						g.perception.Grid[link.Y][link.X] = link
					}
				} else if g.perception.FadeCountdown == 0 {
					switch item := g.perception.Grid[nY][nX].(type) {
					case *Link:
						idx := slices.Index(g.perception.Snakes[item.SnakeId].Links, item)
						if idx > 0 {
							g.biteSnake(item, snake, idx)
						}
					case *Apple:
						g.incScore(snake, model.AppleScore)
						g.perception.Grid[nY][nX] = nil
						g.applePresent = false
					}
				}
			}
			snake.Direction = direction
		}
		if !g.applePresent && rand.IntN(model.NewAppleProbabilityParam) == 0 {
			g.tryToPutAnotherApple()
		}
		g.perception.ElapsedFrames++
	case Scoreboard:
		g.updateScoreboard()
	}
	g.perception.Countdown = int(math.Ceil(float64(g.countdown)/model.Tps)) - 1
	return nil
}

func (g *Game) biteSnake(bittenLink *Link, bitingSnake *Snake, idx int) {
	targetSnake := g.perception.Snakes[bittenLink.SnakeId]
	bittenLink.HealthPercent -= model.HealthReductionPerBite
	bittenLink.Redness = 1
	g.snakeControllers[targetSnake.Id].Vibrate(200 * time.Millisecond)
	if targetSnake != bitingSnake {
		g.incScore(bitingSnake, model.BitLinkScore)
	}
	if bittenLink.HealthPercent <= 0 {
		if targetSnake != bitingSnake {
			g.incScore(bitingSnake, (len(targetSnake.Links)-idx+1)*model.NippedTailLinkBonusMultiplier)
		}
		for i := idx; i < len(targetSnake.Links); i++ {
			link := targetSnake.Links[i]
			g.perception.Grid[link.Y][link.X] = nil
		}
		targetSnake.Links = targetSnake.Links[:idx]
	}
	if bitingSnake.Score >= model.TargetScore {
		g.perception.FadeCountdown = model.GridFadeCountdown
	}
}

func (g *Game) incScore(snake *Snake, delta int) {
	snake.Score += delta
	if snake.Score >= model.TargetScore {
		g.perception.FadeCountdown = model.GridFadeCountdown
	}
}

func (g *Game) updateHeadCount() {
	for _, snake := range g.perception.Snakes {
		snake.Links[0].ChangeRedness(-0.1)
	}
	g.controllers = input.Controllers()
	for _, c := range g.controllers {
		if c.IsAnyJustPressed() {
			snakes := g.perception.Snakes
			snakeCount := len(snakes)
			snakeIdx := slices.IndexFunc(snakes, func(snake *Snake) bool { return g.snakeControllers[snake.Id].Equals(c) })
			if snakeIdx == -1 {
				if snakeCount < model.MaxSnakes {
					for _, snake := range snakes {
						head := snake.Links[0]
						g.perception.Grid[head.Y][head.X] = nil
					}
					g.perception.Snakes = append(g.perception.Snakes, NewSnake(snakeCount, graphics.SnakeColours[snakeCount]))
					g.snakeControllers = append(g.snakeControllers, c)
					g.layoutSnakes()
				}
			} else {
				snakes[snakeIdx].Links[0].Redness = 1
				if c.IsStartJustPressed() && snakeCount > 1 {
					g.perception.Stage = Action
				}
			}
		}
	}
}

func (g *Game) updateScoreboard() {
	for _, snake := range g.perception.Snakes {
		controller := g.snakeControllers[snake.Id]
		if controller.IsStartJustPressed() {
			g.restartPreservingSnakes()
		} else if controller.IsExitJustPressed() {
			os.Exit(0)
		}
		for _, link := range snake.Links {
			link.ChangeRedness(-0.1)
		}
		if controller.IsAnyJustPressed() {
			for _, link := range snake.Links {
				link.Redness = 1
			}
		}
	}
}

func (g *Game) layoutSnakes() {
	delta := 2 * math.Pi / float64(len(g.perception.Snakes))
	alpha := float64(0)
	for _, s := range g.perception.Snakes {
		y := model.GridSize/2 - int(math.Cos(alpha)*model.GridSize/3)
		x := model.GridSize/2 + int(math.Sin(alpha)*model.GridSize/3)
		head := s.Links[0]
		head.X = x
		head.Y = y
		g.perception.Grid[y][x] = head
		alpha += delta
		s.PickInitialDirection()
	}
}

func (g *Game) isAnyButtonPressed(id ebiten.GamepadID) bool {
	buttonPressed := false
	for b := ebiten.StandardGamepadButton(0); b <= ebiten.StandardGamepadButtonMax; b++ {
		if inpututil.IsStandardGamepadButtonJustPressed(id, b) {
			buttonPressed = true
		}
	}
	return buttonPressed
}

func (g *Game) randomUnoccupiedCell() (int, int) {
	x := rand.IntN(model.GridSize)
	y := rand.IntN(model.GridSize)
	for ; y < model.GridSize; y++ {
		for ; x < model.GridSize; x++ {
			if g.perception.Grid[y][x] == nil {
				return x, y
			}
		}
		x = 0
	}
	return -1, -1
}

func (g *Game) tryToPutAnotherApple() {
	x, y := g.randomUnoccupiedCell()
	if x != -1 && y != -1 {
		g.perception.Grid[y][x] = &Apple{X: x, Y: y}
		g.applePresent = true
	}
}

func (g *Game) restartPreservingSnakes() {
	g.perception.Grid = [model.GridSize][model.GridSize]any{}
	g.perception.Stage = Lobby
	g.perception.FadeCountdown = 0
	g.perception.ElapsedFrames = 0
	g.countdown = model.Tps * model.CountdownSeconds
	g.applePresent = false
	g.snakeHeadsRednessGrowth = -1
	g.layoutSnakes()
	for _, snake := range g.perception.Snakes {
		snake.Score = 0
		snake.Links = snake.Links[0:1]
	}
}

func newHeadCoords(s *Snake, direction Direction) (int, int) {
	head := s.Links[0]
	nX := head.X + direction.Dx()
	nY := head.Y + direction.Dy()
	// assuming Dx and Dy can only be -1, 0, 1
	if nX < 0 {
		nX = model.GridSize - 1
	}
	if nY < 0 {
		nY = model.GridSize - 1
	}
	if nX >= model.GridSize {
		nX = 0
	}
	if nY >= model.GridSize {
		nY = 0
	}
	return nX, nY
}
