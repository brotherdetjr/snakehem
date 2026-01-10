package gamestate

import (
	"math"
	"snakehem/internal/interfaces"

	"snakehem/internal/entities"
)

// LobbyState handles player joining before the game starts
type LobbyState struct{}

// Name returns the name of this state
func (s *LobbyState) Name() string {
	return "Lobby"
}

// Update handles the lobby state logic
func (s *LobbyState) Update(ctx *StateContext) (GameState, error) {
	// Fade out redness on all snake heads
	for _, snake := range ctx.Snakes {
		snake.Links[0].ChangeRedness(-0.1, ctx.Config.TpsMultiplier)
	}

	// Process controller inputs
	for _, controller := range ctx.Controllers {
		if controller.IsAnyJustPressed() {
			s.handleControllerPress(ctx, controller)
		}
	}

	// Check if we should transition to Action state
	if s.ShouldTransitionToAction(ctx) {
		return &ActionState{}, nil
	}

	return s, nil
}

// handleControllerPress handles a controller button press in the lobby
func (s *LobbyState) handleControllerPress(ctx *StateContext, controller interfaces.ControllerInput) {
	// Check if this controller is already associated with a snake
	snakeIdx := s.findSnakeByController(ctx.Snakes, controller)

	if snakeIdx == -1 {
		// New player joining
		if len(ctx.Snakes) < ctx.Config.MaxSnakes {
			s.addNewSnake(ctx, controller)
		}
	} else {
		// Existing player - highlight their snake
		ctx.Snakes[snakeIdx].Links[0].Redness = 1
		// State transition is handled in Update() method
	}
}

// findSnakeByController finds a snake controlled by the given controller
func (s *LobbyState) findSnakeByController(snakes []*entities.Snake, controller interfaces.ControllerInput) int {
	for i, snake := range snakes {
		if snake.ControllerID == controller.ID() {
			return i
		}
	}
	return -1
}

// addNewSnake adds a new snake to the game
func (s *LobbyState) addNewSnake(ctx *StateContext, controller interfaces.ControllerInput) {
	// Clear existing snake positions from grid
	for _, snake := range ctx.Snakes {
		head := snake.Links[0]
		ctx.Grid.Clear(head.X, head.Y)
	}

	// Create new snake
	colorIdx := len(ctx.Snakes)
	newSnake := entities.NewSnake(controller.ID(), ctx.Config.SnakeColors[colorIdx])
	ctx.Snakes = append(ctx.Snakes, newSnake)

	// Re-layout all snakes in a circle
	s.layoutSnakes(ctx)
}

// layoutSnakes arranges all snakes in a circle around the grid center
func (s *LobbyState) layoutSnakes(ctx *StateContext) {
	delta := 2 * math.Pi / float64(len(ctx.Snakes))
	alpha := 0.0

	for _, snake := range ctx.Snakes {
		y := ctx.Config.GridSize/2 - int(math.Cos(alpha)*float64(ctx.Config.GridSize)/3)
		x := ctx.Config.GridSize/2 + int(math.Sin(alpha)*float64(ctx.Config.GridSize)/3)

		head := snake.Links[0]
		head.X = x
		head.Y = y
		ctx.Grid.Set(x, y, head)

		alpha += delta

		// Pick initial direction based on position
		snake.PickInitialDirection(ctx.Config.GridSize)
	}
}

// ShouldTransitionToAction checks if we should transition to the action state
func (s *LobbyState) ShouldTransitionToAction(ctx *StateContext) bool {
	if len(ctx.Snakes) < 2 {
		return false
	}

	// Check if any player pressed Start
	for _, controller := range ctx.Controllers {
		if controller.IsStartJustPressed() {
			// Find the snake for this controller
			for _, snake := range ctx.Snakes {
				if snake.ControllerID == controller.ID() {
					return true
				}
			}
		}
	}

	return false
}
