package game

import (
	"os"
	"slices"
	"snakehem/game/common"
	"snakehem/game/local"
	"snakehem/game/shared"
	. "snakehem/game/shared/snake"
	"snakehem/input"
	"snakehem/model"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/rs/zerolog/log"
)

func (g *Game) Update() error {
	start := time.Now()
	defer func() {
		g.unshadedContent.RecordUpdateTimeAndTps(start)
	}()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		log.Info().Msg("Exiting game")
		os.Exit(0)
	}
	g.localContent.Update(&common.Context{Tick: ebiten.Tick()})
	g.unshadedContent.Update()
	switch g.sharedContent.Stage {
	case shared.Lobby:
		g.updateHeadCount()
	case shared.Action:
		if g.sharedContent.GetCountdownSeconds() >= 0 {
			g.sharedContent.Countdown--
		}
		if g.sharedContent.FadeCountdown > 0 {
			g.sharedContent.FadeCountdown--
			if g.sharedContent.FadeCountdown == 0 {
				g.sharedContent.SwitchToScoreboardStage()
				break
			}
		}
		for _, snake := range g.sharedContent.Snakes {
			head := snake.Links[0]
			if g.sharedContent.GetCountdownSeconds() <= 0 {
				var snakeHeadsRednessGrowth float32
				if (g.sharedContent.ActionFrameCount/model.Tps)%2 == 0 {
					snakeHeadsRednessGrowth = -1
				} else {
					snakeHeadsRednessGrowth = 1
				}
				head.ChangeRedness(0.2 * snakeHeadsRednessGrowth)
			} else if g.activeControllers[snake.Id].IsAnyJustPressed() && g.sharedContent.FadeCountdown == 0 {
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
		if g.sharedContent.GetCountdownSeconds() > 0 {
			break
		}
		for _, snake := range g.sharedContent.Snakes {
			direction := snake.Direction
			controller := g.activeControllers[snake.Id]
			if g.sharedContent.FadeCountdown == 0 {
				if controller.IsUpJustPressed() {
					direction = Up
					log.Debug().Int("snakeId", snake.Id).Str("direction", "Up").Msg("New direction")
				} else if controller.IsDownJustPressed() {
					direction = Down
					log.Debug().Int("snakeId", snake.Id).Str("direction", "Down").Msg("New direction")
				} else if controller.IsLeftJustPressed() {
					direction = Left
					log.Debug().Int("snakeId", snake.Id).Str("direction", "Left").Msg("New direction")
				} else if controller.IsRightJustPressed() {
					direction = Right
					log.Debug().Int("snakeId", snake.Id).Str("direction", "Right").Msg("New direction")
				}
			}
			nX, nY := newHeadCoords(snake, direction)
			// not biting self in the neck, preserving same direction if the case
			if len(snake.Links) > 1 && nX == snake.Links[1].X && nY == snake.Links[1].Y {
				direction = snake.Direction
				nX, nY = newHeadCoords(snake, direction)
			}
			if g.sharedContent.ActionFrameCount%model.TpsMultiplier == 0 {
				if g.sharedContent.Grid[nY][nX] == nil {
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
						g.sharedContent.Grid[oldTailY][oldTailX] = nil
					}
					for _, link := range snake.Links {
						g.sharedContent.Grid[link.Y][link.X] = link
					}
					if g.sharedContent.IsAppleHere(nX, nY) {
						g.sharedContent.EatApple(snake)
					}
				} else if g.sharedContent.FadeCountdown == 0 {
					switch item := g.sharedContent.Grid[nY][nX].(type) {
					case *Link:
						idx := slices.Index(g.sharedContent.Snakes[item.SnakeId].Links, item)
						if idx > 0 {
							g.biteSnake(item, snake, idx)
						}
					}
				}
			}
			snake.Direction = direction
		}
		g.sharedContent.TryToPutNewApple()
		g.sharedContent.ActionFrameCount++
	case shared.Scoreboard:
		g.updateScoreboard()
	}
	return nil
}

func (g *Game) biteSnake(bittenLink *Link, bitingSnake *Snake, idx int) {
	targetSnake := g.sharedContent.Snakes[bittenLink.SnakeId]
	bittenLink.HealthPercent -= model.HealthReductionPerBite
	bittenLink.Redness = 1
	g.activeControllers[targetSnake.Id].Vibrate(200 * time.Millisecond)
	if targetSnake != bitingSnake {
		g.sharedContent.IncScore(bitingSnake, model.BitLinkScore)
	}
	if bittenLink.HealthPercent <= 0 {
		if targetSnake != bitingSnake {
			nippedTailLength := len(targetSnake.Links) - idx
			log.Debug().
				Int("bitingSnakeId", bitingSnake.Id).
				Int("targetSnakeId", targetSnake.Id).
				Int("nippedTailLength", nippedTailLength).
				Msg("Nip!")
			g.sharedContent.IncScore(bitingSnake, nippedTailLength*model.NippedTailLinkBonusMultiplier)
		}
		for i := idx; i < len(targetSnake.Links); i++ {
			link := targetSnake.Links[i]
			g.sharedContent.Grid[link.Y][link.X] = nil
		}
		targetSnake.Links = targetSnake.Links[:idx]
	}
}

func (g *Game) updateHeadCount() {
	for _, snake := range g.sharedContent.Snakes {
		snake.Links[0].ChangeRedness(-0.1)
	}
	g.controllers = input.Controllers()
	for _, c := range g.controllers {
		if c.IsAnyJustPressed() {
			snakes := g.sharedContent.Snakes
			snakeCount := len(snakes)
			snakeIdx := slices.IndexFunc(snakes, func(snake *Snake) bool { return g.activeControllers[snake.Id].Equals(c) })
			if snakeIdx == -1 {
				if snakeCount < model.MaxSnakes && g.localContent.GetStage() == local.Off {
					// Start name entry for new player
					g.localContent.SwitchToPlayerNameStage(
						c,
						"Player "+string(rune('0'+(snakeCount+1))),
						common.SnakeColours[snakeCount],
						func(s string) {
							// Submit name and join game
							playerName := strings.TrimSpace(s)
							snakes := g.sharedContent.Snakes
							snakeCount := len(snakes)
							for _, snake := range snakes {
								head := snake.Links[0]
								g.sharedContent.Grid[head.Y][head.X] = nil
							}
							newSnake := NewSnake(snakeCount, playerName, common.SnakeColours[snakeCount])
							g.sharedContent.Snakes = append(g.sharedContent.Snakes, newSnake)
							g.activeControllers = append(g.activeControllers, c)
							g.sharedContent.LayoutSnakes()
							log.Info().Str("name", playerName).Int("id", snakeCount).Msg("Player joined")
						},
					)
				}
			} else {
				snakes[snakeIdx].Links[0].Redness = 1
				if c.IsStartJustPressed() && snakeCount > 1 {
					g.sharedContent.Stage = shared.Action
					log.Info().Int("tagetScore", model.TargetScore).Msg("Action started!")
				}
			}
		}
	}
}

func (g *Game) updateScoreboard() {
	for _, snake := range g.sharedContent.Snakes {
		controller := g.activeControllers[snake.Id]
		if controller.IsStartJustPressed() {
			g.sharedContent.SwitchToLobbyStage()
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

func (g *Game) isAnyButtonPressed(id ebiten.GamepadID) bool {
	buttonPressed := false
	for b := ebiten.StandardGamepadButton(0); b <= ebiten.StandardGamepadButtonMax; b++ {
		if inpututil.IsStandardGamepadButtonJustPressed(id, b) {
			buttonPressed = true
		}
	}
	return buttonPressed
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
