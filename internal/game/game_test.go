package game

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"snakehem/internal/config"
	"snakehem/internal/entities"
	"snakehem/internal/gamestate"
	"snakehem/internal/mocks"
	"snakehem/pkg/ebiten_adapter"
)

// MockCompositeRenderer tracks which rendering methods are called
type MockCompositeRenderer struct {
	drawBackgroundCalled      bool
	drawGridCalled            bool
	drawLobbyUICalled         bool
	drawActionUICalled        bool
	drawScoreboardUICalled    bool
	applyPostProcessingCalled bool
}

func (m *MockCompositeRenderer) DrawBackground(screen *ebiten.Image) {
	m.drawBackgroundCalled = true
}

func (m *MockCompositeRenderer) DrawGrid(screen *ebiten.Image, grid *entities.GameGrid, countdown int) {
	m.drawGridCalled = true
}

func (m *MockCompositeRenderer) DrawLobbyUI(screen *ebiten.Image, snakeCount int) {
	m.drawLobbyUICalled = true
}

func (m *MockCompositeRenderer) DrawActionUI(screen *ebiten.Image, snakes []*entities.Snake, countdown int, elapsedFrames uint64, fadeCountdown int) {
	m.drawActionUICalled = true
}

func (m *MockCompositeRenderer) DrawScoreboardUI(screen *ebiten.Image, snakes []*entities.Snake, elapsedFrames uint64) {
	m.drawScoreboardUICalled = true
}

func (m *MockCompositeRenderer) ApplyPostProcessing(screen *ebiten.Image) {
	m.applyPostProcessingCalled = true
}

func (m *MockCompositeRenderer) Reset() {
	m.drawBackgroundCalled = false
	m.drawGridCalled = false
	m.drawLobbyUICalled = false
	m.drawActionUICalled = false
	m.drawScoreboardUICalled = false
	m.applyPostProcessingCalled = false
}

// Verify interface compliance at compile time
var _ GameRenderer = (*MockCompositeRenderer)(nil)

// TestDraw_LobbyState verifies Bug #1 fix: lobby state should only draw lobby UI
func TestDraw_LobbyState(t *testing.T) {
	cfg := config.DefaultConfig()
	mockRenderer := &MockCompositeRenderer{}
	mockInput := mocks.NewMockInputProvider()
	mockRandom := mocks.NewMockRandomSource([]int{5, 10, 15})

	// Create game in lobby state
	game := &Game{
		config:        cfg,
		inputProvider: mockInput,
		renderer:      mockRenderer,
		random:        mockRandom,
		grid:          entities.NewGameGrid(cfg.GridSize),
		snakes:        make([]*entities.Snake, 0),
		countdown:     cfg.Tps() * cfg.CountdownSeconds,
		elapsedFrames: 0,
		fadeCountdown: 0,
		applePresent:  false,
		currentState:  &gamestate.LobbyState{},
	}

	// Create a mock screen
	screen := ebiten_adapter.NewScreenAdapter(ebiten.NewImage(cfg.GridDimPx(), cfg.GridDimPx()))

	// Call Draw
	game.Draw(screen)

	// Verify: DrawLobbyUI should be called
	if !mockRenderer.drawLobbyUICalled {
		t.Error("DrawLobbyUI should be called in lobby state")
	}

	// Verify: DrawActionUI should NOT be called (Bug #1)
	if mockRenderer.drawActionUICalled {
		t.Error("DrawActionUI should NOT be called in lobby state (Bug #1)")
	}

	// Verify: DrawScoreboardUI should NOT be called
	if mockRenderer.drawScoreboardUICalled {
		t.Error("DrawScoreboardUI should NOT be called in lobby state")
	}

	// Verify: Common rendering methods should be called
	if !mockRenderer.drawBackgroundCalled {
		t.Error("DrawBackground should be called")
	}
	if !mockRenderer.drawGridCalled {
		t.Error("DrawGrid should be called")
	}
	if !mockRenderer.applyPostProcessingCalled {
		t.Error("ApplyPostProcessing should be called")
	}
}

// TestDraw_ActionState verifies that action state draws action UI
func TestDraw_ActionState(t *testing.T) {
	cfg := config.DefaultConfig()
	mockRenderer := &MockCompositeRenderer{}
	mockInput := mocks.NewMockInputProvider()
	mockRandom := mocks.NewMockRandomSource([]int{5, 10, 15})

	// Create game in action state
	game := &Game{
		config:        cfg,
		inputProvider: mockInput,
		renderer:      mockRenderer,
		random:        mockRandom,
		grid:          entities.NewGameGrid(cfg.GridSize),
		snakes:        make([]*entities.Snake, 0),
		countdown:     cfg.Tps() * cfg.CountdownSeconds,
		elapsedFrames: 0,
		fadeCountdown: 0,
		applePresent:  false,
		currentState:  &gamestate.ActionState{},
	}

	// Create a mock screen
	screen := ebiten_adapter.NewScreenAdapter(ebiten.NewImage(cfg.GridDimPx(), cfg.GridDimPx()))

	// Call Draw
	game.Draw(screen)

	// Verify: DrawActionUI should be called
	if !mockRenderer.drawActionUICalled {
		t.Error("DrawActionUI should be called in action state")
	}

	// Verify: DrawLobbyUI should NOT be called
	if mockRenderer.drawLobbyUICalled {
		t.Error("DrawLobbyUI should NOT be called in action state")
	}

	// Verify: DrawScoreboardUI should NOT be called
	if mockRenderer.drawScoreboardUICalled {
		t.Error("DrawScoreboardUI should NOT be called in action state")
	}
}

// TestDraw_ScoreboardState verifies that scoreboard state draws scoreboard UI
func TestDraw_ScoreboardState(t *testing.T) {
	cfg := config.DefaultConfig()
	mockRenderer := &MockCompositeRenderer{}
	mockInput := mocks.NewMockInputProvider()
	mockRandom := mocks.NewMockRandomSource([]int{5, 10, 15})

	// Create game in scoreboard state
	game := &Game{
		config:        cfg,
		inputProvider: mockInput,
		renderer:      mockRenderer,
		random:        mockRandom,
		grid:          entities.NewGameGrid(cfg.GridSize),
		snakes:        make([]*entities.Snake, 0),
		countdown:     cfg.Tps() * cfg.CountdownSeconds,
		elapsedFrames: 0,
		fadeCountdown: 0,
		applePresent:  false,
		currentState:  &gamestate.ScoreboardState{},
	}

	// Create a mock screen
	screen := ebiten_adapter.NewScreenAdapter(ebiten.NewImage(cfg.GridDimPx(), cfg.GridDimPx()))

	// Call Draw
	game.Draw(screen)

	// Verify: DrawScoreboardUI should be called
	if !mockRenderer.drawScoreboardUICalled {
		t.Error("DrawScoreboardUI should be called in scoreboard state")
	}

	// Verify: DrawLobbyUI should NOT be called
	if mockRenderer.drawLobbyUICalled {
		t.Error("DrawLobbyUI should NOT be called in scoreboard state")
	}

	// Verify: DrawActionUI should NOT be called
	if mockRenderer.drawActionUICalled {
		t.Error("DrawActionUI should NOT be called in scoreboard state")
	}
}
