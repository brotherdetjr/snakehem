package game

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"

	"snakehem/internal/config"
	"snakehem/internal/engine"
	"snakehem/internal/entities"
	"snakehem/internal/gamestate"
	"snakehem/internal/interfaces"
	"snakehem/internal/rendering"
)

// ErrUserExit is returned when the user requests to exit the game
var ErrUserExit = errors.New("user requested exit")

// Game orchestrates the game using dependency-injected components
// This is a thin layer that delegates to state management and rendering
type Game struct {
	// Injected dependencies (immutable)
	config        *config.GameConfig
	inputProvider interfaces.InputProvider
	renderer      *rendering.CompositeRenderer
	random        interfaces.RandomSource
	physics       *engine.PhysicsEngine
	scoring       *engine.ScoringEngine

	// Mutable game state
	grid          *entities.GameGrid
	snakes        []*entities.Snake
	countdown     int
	elapsedFrames uint64
	fadeCountdown int
	applePresent  bool
	currentState  gamestate.GameState
}

// Update is called every tick (60 FPS by default)
// Reduced from 113 lines to ~25 lines by delegating to state pattern
func (g *Game) Update() error {
	// Check for global exit (Escape key)
	if g.inputProvider.IsGlobalExitPressed() {
		return ErrUserExit
	}

	// Get all active controllers
	controllers := g.inputProvider.GetControllerInputs()

	// Create state context with all dependencies
	ctx := &gamestate.StateContext{
		Snakes:        g.snakes,
		Grid:          g.grid,
		Controllers:   controllers,
		Physics:       g.physics,
		Scoring:       g.scoring,
		Random:        g.random,
		Config:        g.config,
		Countdown:     &g.countdown,
		ElapsedFrames: &g.elapsedFrames,
		FadeCountdown: &g.fadeCountdown,
		ApplePresent:  &g.applePresent,
	}

	// Delegate to current state
	nextState, err := g.currentState.Update(ctx)
	if err != nil {
		return err
	}

	// Update current state if transition occurred
	if nextState != g.currentState {
		g.currentState = nextState
	}

	// Increment elapsed frames
	g.elapsedFrames++

	return nil
}

// Draw renders the game to the screen
// Delegates to the composite renderer based on current state
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background
	g.renderer.DrawBackground(screen)

	// Draw grid (snakes and apples)
	g.renderer.DrawGrid(screen, g.grid, g.countdown)

	// Draw state-specific UI
	switch state := g.currentState.(type) {
	case *gamestate.LobbyState:
		g.renderer.DrawLobbyUI(screen, len(g.snakes))
		g.renderer.DrawActionUI(screen, g.snakes, g.countdown, g.elapsedFrames, 0)

	case *gamestate.ActionState:
		g.renderer.DrawActionUI(screen, g.snakes, g.countdown, g.elapsedFrames, g.fadeCountdown)

	case *gamestate.ScoreboardState:
		g.renderer.DrawScoreboardUI(screen, g.snakes, g.elapsedFrames)

	default:
		// Unknown state - shouldn't happen
		_ = state
	}

	// Apply post-processing (CRT shader)
	g.renderer.ApplyPostProcessing(screen)
}

// Layout returns the game's logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.config.GridDimPx(), g.config.GridDimPx()
}
