package rendering

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"snakehem/internal/config"
	"snakehem/internal/entities"
)

// TestDrawLobbyUI_WithSnakes verifies that DrawLobbyUI accepts snakes parameter
func TestDrawLobbyUI_WithSnakes(t *testing.T) {
	cfg := config.DefaultConfig()

	// Create components
	fontManager := NewFontManager()
	snakeRenderer := NewSnakeRenderer(cfg)
	uiRenderer := NewUIRenderer(cfg, fontManager)
	postProcessor, err := NewPostProcessor()
	if err != nil {
		t.Fatalf("Failed to create post processor: %v", err)
	}

	renderer := NewCompositeRenderer(
		cfg,
		fontManager,
		snakeRenderer,
		uiRenderer,
		postProcessor,
	)

	// Create test data
	snakes := []*entities.Snake{
		entities.NewSnake("player1", color.RGBA{R: 255, G: 255, B: 255, A: 255}),
	}

	// Create screen
	screen := ebiten.NewImage(cfg.GridDimPx(), cfg.GridDimPx())

	// This should not panic
	renderer.DrawLobbyUI(screen, snakes, 0, cfg.Tps()*cfg.CountdownSeconds)

	// If we get here, the method signature is correct and it executed
	t.Log("DrawLobbyUI executed successfully with snakes parameter")
}

// TestDrawScores_InLobbyState verifies DrawScores is called with correct parameters
func TestDrawScores_InLobbyState(t *testing.T) {
	cfg := config.DefaultConfig()

	fontManager := NewFontManager()
	uiRenderer := NewUIRenderer(cfg, fontManager)

	// Create test snake with score
	snake := entities.NewSnake("player1", color.RGBA{R: 200, G: 200, B: 200, A: 255})
	snake.Score = 0
	snake.Links[0].Redness = 1.0 // Full red

	snakes := []*entities.Snake{snake}

	screen := ebiten.NewImage(cfg.GridDimPx(), cfg.GridDimPx())

	// This should not panic - isActionState=false for lobby
	uiRenderer.DrawScores(screen, snakes, 0, false, cfg.Tps()*cfg.CountdownSeconds)

	t.Log("DrawScores executed successfully for lobby state")
}
