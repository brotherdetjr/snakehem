package gamestate

import (
	"snakehem/internal/entities"
)

// ScoreboardState handles the post-game scoreboard
type ScoreboardState struct{}

// Name returns the name of this state
func (s *ScoreboardState) Name() string {
	return "Scoreboard"
}

// Update handles the scoreboard state logic
func (s *ScoreboardState) Update(ctx *StateContext) (GameState, error) {
	// Process controller inputs
	for _, controller := range ctx.Controllers {
		// Find the snake for this controller
		var controllerSnake *entities.Snake
		for _, snake := range ctx.Snakes {
			if snake.ControllerID == controller.ID() {
				controllerSnake = snake
				break
			}
		}

		if controllerSnake == nil {
			continue
		}

		// Check for restart
		if controller.IsStartJustPressed() {
			s.restartGame(ctx)
			return &LobbyState{}, nil
		}

		// Update redness
		for _, link := range controllerSnake.Links {
			link.ChangeRedness(-0.1, ctx.Config.TpsMultiplier)
		}

		// Highlight on any button press
		if controller.IsAnyJustPressed() {
			controllerSnake.Links[0].Redness = 1
		}
	}

	return s, nil
}

// restartGame resets the game state while preserving snakes
func (s *ScoreboardState) restartGame(ctx *StateContext) {
	// Clear the grid
	ctx.Grid.ClearAll()

	// Reset game state variables
	*ctx.Countdown = ctx.Config.Tps() * ctx.Config.CountdownSeconds
	*ctx.ElapsedFrames = 0
	*ctx.FadeCountdown = 0
	*ctx.ApplePresent = false

	// Reset each snake
	for _, snake := range ctx.Snakes {
		snake.Score = 0
		snake.Links = snake.Links[0:1]
		snake.HeadRednessGrowth = -1
	}

	// Re-layout snakes (reusing lobby state's layout logic)
	lobby := &LobbyState{}
	lobby.layoutSnakes(ctx)
}
