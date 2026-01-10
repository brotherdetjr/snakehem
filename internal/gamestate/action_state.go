package gamestate

import (
	"time"

	"github.com/rs/zerolog/log"

	"snakehem/direction"
	"snakehem/internal/entities"
)

// ActionState handles the main gameplay
type ActionState struct{}

// Name returns the name of this state
func (s *ActionState) Name() string {
	return "Action"
}

// Update handles the action state logic
func (s *ActionState) Update(ctx *StateContext) (GameState, error) {
	// Update countdown
	if *ctx.Countdown > 0 {
		*ctx.Countdown--
	}

	// Update fade countdown
	if *ctx.FadeCountdown > 0 {
		*ctx.FadeCountdown--
		if *ctx.FadeCountdown == 0 {
			// Transition to scoreboard
			return &ScoreboardState{}, nil
		}
	}

	// Update snake redness
	s.updateSnakeRedness(ctx)

	// Skip movement during countdown
	if *ctx.Countdown > ctx.Config.Tps() {
		return s, nil
	}

	// Process snake movements and collisions
	if *ctx.ElapsedFrames%uint64(ctx.Config.TpsMultiplier) == 0 {
		s.processSnakeMovements(ctx)
	}

	// Spawn apples randomly
	if !*ctx.ApplePresent && ctx.Random.IntN(ctx.Config.NewAppleProbabilityParam) == 0 {
		s.trySpawnApple(ctx)
	}

	*ctx.ElapsedFrames++

	return s, nil
}

// updateSnakeRedness updates the visual redness effect on snakes
func (s *ActionState) updateSnakeRedness(ctx *StateContext) {
	for _, snake := range ctx.Snakes {
		head := snake.Links[0]

		if *ctx.Countdown <= ctx.Config.Tps() {
			// Pulsing effect during gameplay
			head.ChangeRedness(0.2*snake.HeadRednessGrowth, ctx.Config.TpsMultiplier)
			if head.Redness >= 1 || head.Redness <= 0 {
				snake.HeadRednessGrowth = -snake.HeadRednessGrowth
			}
		} else {
			// Pre-game countdown
			controller, exists := ctx.Controllers[snake.ControllerID]
			if exists && controller.IsAnyJustPressed() && *ctx.FadeCountdown == 0 {
				head.Redness = 1
			} else {
				head.ChangeRedness(-0.1, ctx.Config.TpsMultiplier)
			}
		}

		// Fade non-head links
		for _, link := range snake.Links {
			if link != head {
				link.ChangeRedness(-0.1, ctx.Config.TpsMultiplier)
			}
		}
	}
}

// processSnakeMovements handles input, movement, and collisions for all snakes
func (s *ActionState) processSnakeMovements(ctx *StateContext) {
	for _, snake := range ctx.Snakes {
		// Determine new direction from input
		newDirection := s.getDirectionFromInput(ctx, snake)

		// Calculate new head position
		nX, nY := ctx.Physics.CalculateNewHeadPosition(snake, newDirection)

		// Check for self-neck bite (don't allow moving into own neck)
		if len(snake.Links) > 1 && nX == snake.Links[1].X && nY == snake.Links[1].Y {
			// Keep current direction instead
			newDirection = snake.Direction
			nX, nY = ctx.Physics.CalculateNewHeadPosition(snake, newDirection)
		}

		// Check what's at the new position
		targetCell := ctx.Grid.Get(nX, nY)

		if targetCell == nil {
			// Empty cell - move snake
			s.moveSnake(ctx, snake, nX, nY)
		} else if *ctx.FadeCountdown == 0 {
			// Collision detected
			s.handleCollision(ctx, snake, nX, nY, targetCell)
		}

		// Update direction
		snake.Direction = newDirection
	}
}

// getDirectionFromInput determines the new direction based on controller input
func (s *ActionState) getDirectionFromInput(ctx *StateContext, snake *entities.Snake) direction.Direction {
	if *ctx.FadeCountdown != 0 {
		return snake.Direction
	}

	controller, exists := ctx.Controllers[snake.ControllerID]
	if !exists {
		return snake.Direction
	}

	if controller.IsUpJustPressed() {
		log.Info().Str("snakeId", snake.ControllerID).Str("direction", "Up").Msg("New direction")
		return direction.Up
	} else if controller.IsDownJustPressed() {
		log.Info().Str("snakeId", snake.ControllerID).Str("direction", "Down").Msg("New direction")
		return direction.Down
	} else if controller.IsLeftJustPressed() {
		log.Info().Str("snakeId", snake.ControllerID).Str("direction", "Left").Msg("New direction")
		return direction.Left
	} else if controller.IsRightJustPressed() {
		log.Info().Str("snakeId", snake.ControllerID).Str("direction", "Right").Msg("New direction")
		return direction.Right
	}

	return snake.Direction
}

// moveSnake moves a snake to a new position
func (s *ActionState) moveSnake(ctx *StateContext, snake *entities.Snake, nX, nY int) {
	tail := snake.Links[len(snake.Links)-1]
	oldTailX := tail.X
	oldTailY := tail.Y

	// Move all links forward
	for i := len(snake.Links) - 1; i > 0; i-- {
		link := snake.Links[i]
		prevLink := snake.Links[i-1]
		link.X = prevLink.X
		link.Y = prevLink.Y
	}

	// Move head to new position
	snake.Links[0].X = nX
	snake.Links[0].Y = nY

	// Grow snake if not at target length
	if len(snake.Links) < ctx.Config.SnakeTargetLength {
		snake.Links = append(snake.Links, &entities.Link{
			HealthPercent: 100,
			Snake:         snake,
			X:             oldTailX,
			Y:             oldTailY,
			Redness:       0,
		})
	} else {
		// Clear old tail position
		ctx.Grid.Clear(oldTailY, oldTailX)
	}

	// Update grid with new positions
	for _, link := range snake.Links {
		ctx.Grid.Set(link.X, link.Y, link)
	}
}

// handleCollision handles collision with another object
func (s *ActionState) handleCollision(ctx *StateContext, snake *entities.Snake, x, y int, targetCell interface{}) {
	switch item := targetCell.(type) {
	case *entities.Link:
		// Find the index of this link in its snake
		idx := s.findLinkIndex(item)
		if idx > 0 {
			// Bite the link (not the head)
			s.biteSnake(ctx, item, snake, idx)
		}
	case *entities.Apple:
		// Eat the apple
		points := ctx.Scoring.ProcessApple(snake, ctx.FadeCountdown)
		snake.Score += points
		ctx.Grid.Clear(x, y)
		*ctx.ApplePresent = false
	}
}

// findLinkIndex finds the index of a link in its snake
func (s *ActionState) findLinkIndex(link *entities.Link) int {
	for i, l := range link.Snake.Links {
		if l == link {
			return i
		}
	}
	return -1
}

// biteSnake handles one snake biting another
func (s *ActionState) biteSnake(ctx *StateContext, bittenLink *entities.Link, bitingSnake *entities.Snake, idx int) {
	targetSnake := bittenLink.Snake

	// Reduce health
	bittenLink.HealthPercent -= ctx.Config.HealthReductionPerBite
	bittenLink.Redness = 1

	// Vibrate the victim's controller
	if controller, exists := ctx.Controllers[targetSnake.ControllerID]; exists {
		controller.Vibrate(200 * time.Millisecond)
	}

	// Award points
	points := ctx.Scoring.ProcessBite(bitingSnake, targetSnake, idx, ctx.Grid, ctx.FadeCountdown)
	bitingSnake.Score += points

	// Remove tail if health depleted
	if bittenLink.HealthPercent <= 0 {
		s.removeTail(ctx, targetSnake, idx)
	}
}

// removeTail removes the tail portion of a snake starting from the given index
func (s *ActionState) removeTail(ctx *StateContext, snake *entities.Snake, fromIndex int) {
	// Clear grid cells
	for i := fromIndex; i < len(snake.Links); i++ {
		link := snake.Links[i]
		ctx.Grid.Clear(link.X, link.Y)
	}

	// Truncate links
	snake.Links = snake.Links[:fromIndex]
}

// trySpawnApple attempts to spawn an apple at a random empty location
func (s *ActionState) trySpawnApple(ctx *StateContext) {
	x, y, found := ctx.Grid.FindRandomEmpty(ctx.Random)
	if found {
		ctx.Grid.Set(x, y, &entities.Apple{X: x, Y: y})
		*ctx.ApplePresent = true
	}
}
