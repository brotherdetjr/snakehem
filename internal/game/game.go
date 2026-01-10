package game

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"

	"snakehem/internal/config"
	"snakehem/internal/engine"
	"snakehem/internal/entities"
	"snakehem/internal/gamestate"
	"snakehem/internal/interfaces"
	"snakehem/pkg/ebiten_adapter"
)

// ErrUserExit is returned when the user requests to exit the game
var ErrUserExit = errors.New("user requested exit")

// GameRenderer defines the interface for the composite game renderer
// This allows for testing with mock implementations
type GameRenderer interface {
	DrawBackground(screen *ebiten.Image)
	DrawGrid(screen *ebiten.Image, grid *entities.GameGrid, countdown int)
	DrawLobbyUI(screen *ebiten.Image, snakeCount int)
	DrawActionUI(screen *ebiten.Image, snakes []*entities.Snake, countdown int, elapsedFrames uint64, fadeCountdown int)
	DrawScoreboardUI(screen *ebiten.Image, snakes []*entities.Snake, elapsedFrames uint64)
	ApplyPostProcessing(screen *ebiten.Image)
}

// Game orchestrates the game using dependency-injected components
// This is a thin layer that delegates to state management and rendering
type Game struct {
	// Injected dependencies (immutable)
	config        *config.GameConfig
	inputProvider interfaces.InputProvider
	renderer      GameRenderer
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

	// Get all active controllers and convert to map
	controllerSlice := g.inputProvider.GetControllerInputs()
	controllers := make(map[string]interfaces.ControllerInput, len(controllerSlice))
	for _, ctrl := range controllerSlice {
		controllers[ctrl.ID()] = ctrl
	}

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
func (g *Game) Draw(screen interfaces.Screen) {
	// Extract *ebiten.Image from interfaces.Screen
	// Renderers use Ebiten-specific functions (like vector drawing) so need concrete type
	var ebitenScreen *ebiten.Image
	if adapter, ok := screen.(*ebiten_adapter.ScreenAdapter); ok {
		ebitenScreen = adapter.EbitenImage()
	} else {
		// Fallback: this shouldn't happen in production
		return
	}

	// Draw background
	g.renderer.DrawBackground(ebitenScreen)

	// Draw grid (snakes and apples)
	g.renderer.DrawGrid(ebitenScreen, g.grid, g.countdown)

	// Draw state-specific UI
	switch state := g.currentState.(type) {
	case *gamestate.LobbyState:
		g.renderer.DrawLobbyUI(ebitenScreen, len(g.snakes))

	case *gamestate.ActionState:
		g.renderer.DrawActionUI(ebitenScreen, g.snakes, g.countdown, g.elapsedFrames, g.fadeCountdown)

	case *gamestate.ScoreboardState:
		g.renderer.DrawScoreboardUI(ebitenScreen, g.snakes, g.elapsedFrames)

	default:
		// Unknown state - shouldn't happen
		_ = state
	}

	// Apply post-processing (CRT shader)
	g.renderer.ApplyPostProcessing(ebitenScreen)
}

// Layout returns the game's logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.config.GridDimPx(), g.config.GridDimPx()
}
