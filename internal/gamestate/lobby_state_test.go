package gamestate

import (
	"image/color"
	"testing"

	"snakehem/internal/config"
	"snakehem/internal/entities"
	"snakehem/internal/interfaces"
	"snakehem/internal/mocks"
)

func TestLobbyState_ShouldTransitionToAction(t *testing.T) {
	cfg := &config.GameConfig{
		GridSize:    10,
		CellDimPx:   16,
		TargetScore: 100,
		MaxSnakes:   4,
	}

	tests := []struct {
		name             string
		snakeCount       int
		controllerIDs    []string
		startPressedOnID string
		wantTransition   bool
	}{
		{
			name:             "not enough players",
			snakeCount:       1,
			controllerIDs:    []string{"player1"},
			startPressedOnID: "player1",
			wantTransition:   false,
		},
		{
			name:             "two players, start pressed",
			snakeCount:       2,
			controllerIDs:    []string{"player1", "player2"},
			startPressedOnID: "player1",
			wantTransition:   true,
		},
		{
			name:             "two players, no start pressed",
			snakeCount:       2,
			controllerIDs:    []string{"player1", "player2"},
			startPressedOnID: "",
			wantTransition:   false,
		},
		{
			name:             "three players, start pressed by second player",
			snakeCount:       3,
			controllerIDs:    []string{"player1", "player2", "player3"},
			startPressedOnID: "player2",
			wantTransition:   true,
		},
		{
			name:             "two players, start pressed by non-player",
			snakeCount:       2,
			controllerIDs:    []string{"player1", "player2"},
			startPressedOnID: "spectator",
			wantTransition:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create snakes
			snakes := make([]*entities.Snake, tt.snakeCount)
			for i := 0; i < tt.snakeCount; i++ {
				snakes[i] = entities.NewSnake(tt.controllerIDs[i], color.RGBA{R: 255, G: 0, B: 0, A: 255})
			}

			// Create controllers map with proper interface type
			controllers := make(map[string]interfaces.ControllerInput)
			for _, id := range tt.controllerIDs {
				ctrl := mocks.NewMockControllerInput(id)
				if id == tt.startPressedOnID {
					ctrl.SetStartPressed(true)
				}
				controllers[id] = ctrl
			}

			// Add spectator controller if needed
			if tt.startPressedOnID == "spectator" {
				ctrl := mocks.NewMockControllerInput("spectator")
				ctrl.SetStartPressed(true)
				controllers["spectator"] = ctrl
			}

			// Create context
			ctx := &StateContext{
				Snakes:      snakes,
				Controllers: controllers,
				Config:      cfg,
			}

			state := &LobbyState{}
			got := state.ShouldTransitionToAction(ctx)

			if got != tt.wantTransition {
				t.Errorf("ShouldTransitionToAction() = %v, want %v", got, tt.wantTransition)
			}
		})
	}
}

func TestLobbyState_FindSnakeByController(t *testing.T) {
	cfg := &config.GameConfig{
		GridSize:  10,
		CellDimPx: 16,
	}

	tests := []struct {
		name      string
		snakeIDs  []string
		searchID  string
		wantIndex int
	}{
		{
			name:      "find first snake",
			snakeIDs:  []string{"player1", "player2", "player3"},
			searchID:  "player1",
			wantIndex: 0,
		},
		{
			name:      "find middle snake",
			snakeIDs:  []string{"player1", "player2", "player3"},
			searchID:  "player2",
			wantIndex: 1,
		},
		{
			name:      "find last snake",
			snakeIDs:  []string{"player1", "player2", "player3"},
			searchID:  "player3",
			wantIndex: 2,
		},
		{
			name:      "controller not found",
			snakeIDs:  []string{"player1", "player2"},
			searchID:  "player3",
			wantIndex: -1,
		},
		{
			name:      "empty snake list",
			snakeIDs:  []string{},
			searchID:  "player1",
			wantIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = cfg // suppress unused warning

			// Create snakes
			snakes := make([]*entities.Snake, len(tt.snakeIDs))
			for i, id := range tt.snakeIDs {
				snakes[i] = entities.NewSnake(id, color.RGBA{R: 255, G: 0, B: 0, A: 255})
			}

			// Create search controller
			searchCtrl := mocks.NewMockControllerInput(tt.searchID)

			state := &LobbyState{}
			got := state.findSnakeByController(snakes, searchCtrl)

			if got != tt.wantIndex {
				t.Errorf("findSnakeByController() = %d, want %d", got, tt.wantIndex)
			}
		})
	}
}

func TestLobbyState_Update_FadeRedness(t *testing.T) {
	cfg := &config.GameConfig{
		GridSize:      10,
		CellDimPx:     16,
		TpsMultiplier: 1.0,
		TargetScore:   100,
	}

	// Create a snake with high redness
	snake := entities.NewSnake("player1", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	snake.Links[0].Redness = 1.0

	grid := entities.NewGameGrid(cfg.GridSize)
	grid.Set(snake.Links[0].X, snake.Links[0].Y, snake.Links[0])

	ctx := &StateContext{
		Snakes:      []*entities.Snake{snake},
		Grid:        grid,
		Controllers: make(map[string]interfaces.ControllerInput),
		Config:      cfg,
	}

	state := &LobbyState{}

	// Update several times
	for i := 0; i < 5; i++ {
		nextState, err := state.Update(ctx)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}
		if nextState != state {
			t.Fatalf("Update() should stay in lobby state")
		}
	}

	// Redness should have decreased
	if snake.Links[0].Redness >= 1.0 {
		t.Errorf("Redness should have decreased, got %f", snake.Links[0].Redness)
	}
	if snake.Links[0].Redness < 0 {
		t.Errorf("Redness should not go below 0, got %f", snake.Links[0].Redness)
	}
}
