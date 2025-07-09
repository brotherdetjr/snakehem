package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/rs/zerolog/log"
	"math"
	"math/rand/v2"
	"os"
	"slices"
	. "snakehem/apple"
	"snakehem/controller"
	"snakehem/controller/gamepad"
	"snakehem/controller/keyboard"
	. "snakehem/direction"
	. "snakehem/snake"
	. "snakehem/state"
	"time"
)

func (g *Game) Update() error {
	if keyboard.Instance.IsExitJustPressed() {
		os.Exit(0)
	}
	switch g.state {
	case Lobby:
		g.updateHeadCount()
	case Action:
		if g.countdown > 0 {
			g.countdown--
		}
		if g.fadeCountdown > 0 {
			g.fadeCountdown--
			if g.fadeCountdown == 0 {
				g.state = Scoreboard
				break
			}
		}
		for _, snake := range g.snakes {
			head := snake.Links[0]
			if g.countdown <= tps {
				head.ChangeRedness(0.2 * snake.HeadRednessGrowth)
				if head.Redness >= 1 || head.Redness <= 0 {
					snake.HeadRednessGrowth = -snake.HeadRednessGrowth
				}
			} else if snake.Controller.IsAnyJustPressed() && g.fadeCountdown == 0 {
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
		if g.countdown > tps {
			break
		}
		if g.countdown == tps {
			for _, snake := range g.snakes {
				snake.Direction = getInitialDirection(snake)
			}
		}
		for _, snake := range g.snakes {
			direction := snake.Direction
			if g.fadeCountdown == 0 {
				if snake.Controller.IsUpJustPressed() {
					direction = Up
					log.Info().Any("snakeId", snake.Controller).Str("direction", "Up").Msg("New direction")
				} else if snake.Controller.IsDownJustPressed() {
					direction = Down
					log.Info().Any("snakeId", snake.Controller).Str("direction", "Down").Msg("New direction")
				} else if snake.Controller.IsLeftJustPressed() {
					direction = Left
					log.Info().Any("snakeId", snake.Controller).Str("direction", "Left").Msg("New direction")
				} else if snake.Controller.IsRightJustPressed() {
					direction = Right
					log.Info().Any("snakeId", snake.Controller).Str("direction", "Right").Msg("New direction")
				}
			}
			nX, nY := newHeadCoords(snake, direction)
			// not biting self in the neck, preserving same direction if the case
			if len(snake.Links) > 1 && nX == snake.Links[1].X && nY == snake.Links[1].Y {
				direction = snake.Direction
				nX, nY = newHeadCoords(snake, direction)
			}
			if g.grid[nY][nX] == nil {
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
				if len(snake.Links) < snakeTargetLength {
					snake.Links = append(snake.Links, &Link{
						HealthPercent: 100,
						Snake:         snake,
						X:             oldTailX,
						Y:             oldTailY,
						Redness:       0,
					})
				} else {
					g.grid[oldTailY][oldTailX] = nil
				}
				for _, link := range snake.Links {
					g.grid[link.Y][link.X] = link
				}
			} else if g.fadeCountdown == 0 {
				switch item := g.grid[nY][nX].(type) {
				case *Link:
					idx := slices.Index(item.Snake.Links, item)
					if idx > 0 {
						g.biteSnake(item, snake, idx)
					}
				case *Apple:
					g.incScore(snake, appleScore)
					g.grid[nY][nX] = nil
					g.applePresent = false
				}
			}
			snake.Direction = direction
		}
		if !g.applePresent && rand.IntN(newAppleProbabilityParam) == 0 {
			g.tryToPutAnotherApple()
		}
		g.elapsedFrames++
	case Scoreboard:
		g.updateScoreboard()
	}
	return nil
}

func (g *Game) biteSnake(bittenLink *Link, bitingSnake *Snake, idx int) {
	targetSnake := bittenLink.Snake
	bittenLink.HealthPercent -= healthReductionPerBite
	bittenLink.Redness = 1
	targetSnake.Controller.Vibrate(200 * time.Millisecond)
	if targetSnake != bitingSnake {
		g.incScore(bitingSnake, bitLinkScore)
	}
	if bittenLink.HealthPercent <= 0 {
		if targetSnake != bitingSnake {
			g.incScore(bitingSnake, (len(targetSnake.Links)-idx+1)*nippedTailLinkBonusMultiplier)
		}
		for i := idx; i < len(targetSnake.Links); i++ {
			link := targetSnake.Links[i]
			g.grid[link.Y][link.X] = nil
		}
		targetSnake.Links = targetSnake.Links[:idx]
	}
	if bitingSnake.Score >= targetScore {
		g.fadeCountdown = gridFadeCountdown
	}
}

func (g *Game) incScore(snake *Snake, delta int) {
	snake.Score += delta
	if snake.Score >= targetScore {
		g.fadeCountdown = gridFadeCountdown
	}
}

func appendControllers(controllers []controller.Controller) []controller.Controller {
	var result []controller.Controller = nil
	if !slices.Contains(controllers, keyboard.Instance) {
		result = append(result, keyboard.Instance)
	}
	for _, g := range ebiten.AppendGamepadIDs(nil) {
		contains := false
		var gamepadAsController controller.Controller = gamepad.NewGamepad(g)
		for _, c := range controllers {
			if c.Equals(gamepadAsController) {
				contains = true
				break
			}
		}
		if !contains && ebiten.IsStandardGamepadLayoutAvailable(g) {
			result = append(result, gamepadAsController)
		}
	}
	return result
}

func (g *Game) updateHeadCount() {
	for _, snake := range g.snakes {
		snake.Links[0].ChangeRedness(-0.1)
	}
	g.controllers = appendControllers(g.controllers)
	for _, c := range g.controllers {
		if c.IsAnyJustPressed() {
			snakeIdx := slices.IndexFunc(g.snakes, func(snake *Snake) bool { return snake.Controller.Equals(c) })
			if snakeIdx == -1 {
				if len(g.snakes) < maxSnakes {
					for _, snake := range g.snakes {
						head := snake.Links[0]
						g.grid[head.Y][head.X] = nil
					}
					g.snakes = append(g.snakes, NewSnake(c, snakeColours[len(g.snakes)]))
					g.layoutSnakes()
				}
			} else {
				g.snakes[snakeIdx].Links[0].Redness = 1
				if c.IsStartJustPressed() && len(g.snakes) > 1 {
					g.state = Action
				}
			}
		}
	}
}

func (g *Game) updateScoreboard() {
	for _, snake := range g.snakes {
		if snake.Controller.IsStartJustPressed() {
			g.restartPreservingSnakes()
		} else if snake.Controller.IsExitJustPressed() {
			os.Exit(0)
		}
		for _, link := range snake.Links {
			link.ChangeRedness(-0.1)
		}
		if snake.Controller.IsAnyJustPressed() {
			snake.Links[0].Redness = 1
		}
	}
}

func (g *Game) layoutSnakes() {
	delta := 2 * math.Pi / float64(len(g.snakes))
	alpha := float64(0)
	for _, s := range g.snakes {
		y := gridSize/2 - int(math.Cos(alpha)*gridSize/3)
		x := gridSize/2 + int(math.Sin(alpha)*gridSize/3)
		head := s.Links[0]
		head.X = x
		head.Y = y
		g.grid[y][x] = head
		alpha += delta
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
	x := rand.IntN(gridSize)
	y := rand.IntN(gridSize)
	for ; y < gridSize; y++ {
		for ; x < gridSize; x++ {
			if g.grid[y][x] == nil {
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
		g.grid[y][x] = &Apple{X: x, Y: y}
		g.applePresent = true
	}
}

func (g *Game) restartPreservingSnakes() {
	g.grid = [gridSize][gridSize]any{}
	g.state = Lobby
	g.countdown = tps * countdownSeconds
	g.elapsedFrames = 0
	g.fadeCountdown = 0
	g.applePresent = false
	g.layoutSnakes()
	for _, snake := range g.snakes {
		snake.Score = 0
		snake.Links = snake.Links[0:1]
		snake.Direction = getInitialDirection(snake)
		snake.HeadRednessGrowth = -1
	}
}

func getInitialDirection(s *Snake) Direction {
	head := s.Links[0]
	x := head.X
	y := head.Y
	midPoint := gridSize/2 + 1
	if math.Abs(float64(midPoint-x)) > math.Abs(float64(midPoint-y)) {
		if midPoint < x {
			return Left
		} else {
			return Right
		}
	} else {
		if midPoint < y {
			return Up
		} else {
			return Down
		}
	}
}

func newHeadCoords(s *Snake, direction Direction) (int, int) {
	head := s.Links[0]
	nX := head.X + direction.Dx()
	nY := head.Y + direction.Dy()
	// assuming Dx and Dy can only be -1, 0, 1
	if nX < 0 {
		nX = gridSize - 1
	}
	if nY < 0 {
		nY = gridSize - 1
	}
	if nX >= gridSize {
		nX = 0
	}
	if nY >= gridSize {
		nY = 0
	}
	return nX, nY
}
