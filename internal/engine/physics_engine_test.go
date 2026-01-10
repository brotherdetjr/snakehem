package engine

import (
	"testing"

	"snakehem/direction"
	"snakehem/internal/config"
	"snakehem/internal/entities"
)

func TestPhysicsEngine_CalculateNewHeadPosition(t *testing.T) {
	cfg := &config.GameConfig{
		GridSize:    10,
		CellDimPx:   16,
		TargetScore: 100,
	}
	engine := NewPhysicsEngine(cfg)

	tests := []struct {
		name      string
		headX     int
		headY     int
		direction direction.Direction
		wantX     int
		wantY     int
	}{
		{
			name:      "move right from center",
			headX:     5,
			headY:     5,
			direction: direction.Right,
			wantX:     6,
			wantY:     5,
		},
		{
			name:      "move left from center",
			headX:     5,
			headY:     5,
			direction: direction.Left,
			wantX:     4,
			wantY:     5,
		},
		{
			name:      "move up from center",
			headX:     5,
			headY:     5,
			direction: direction.Up,
			wantX:     5,
			wantY:     4,
		},
		{
			name:      "move down from center",
			headX:     5,
			headY:     5,
			direction: direction.Down,
			wantX:     5,
			wantY:     6,
		},
		{
			name:      "wrap around right edge",
			headX:     9,
			headY:     5,
			direction: direction.Right,
			wantX:     0,
			wantY:     5,
		},
		{
			name:      "wrap around left edge",
			headX:     0,
			headY:     5,
			direction: direction.Left,
			wantX:     9,
			wantY:     5,
		},
		{
			name:      "wrap around top edge",
			headX:     5,
			headY:     0,
			direction: direction.Up,
			wantX:     5,
			wantY:     9,
		},
		{
			name:      "wrap around bottom edge",
			headX:     5,
			headY:     9,
			direction: direction.Down,
			wantX:     5,
			wantY:     0,
		},
		{
			name:      "wrap around top-left corner",
			headX:     0,
			headY:     0,
			direction: direction.Left,
			wantX:     9,
			wantY:     0,
		},
		{
			name:      "wrap around bottom-right corner",
			headX:     9,
			headY:     9,
			direction: direction.Down,
			wantX:     9,
			wantY:     0,
		},
		{
			name:      "no movement with None direction",
			headX:     5,
			headY:     5,
			direction: direction.None,
			wantX:     5,
			wantY:     5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a snake with head at the specified position
			snake := &entities.Snake{
				Links: []*entities.Link{
					{X: tt.headX, Y: tt.headY},
				},
				Direction: tt.direction,
			}

			gotX, gotY := engine.CalculateNewHeadPosition(snake, tt.direction)

			if gotX != tt.wantX || gotY != tt.wantY {
				t.Errorf("CalculateNewHeadPosition() = (%d, %d), want (%d, %d)",
					gotX, gotY, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestPhysicsEngine_CalculateNewHeadPosition_DifferentGridSizes(t *testing.T) {
	tests := []struct {
		name      string
		gridSize  int
		headX     int
		headY     int
		direction direction.Direction
		wantX     int
		wantY     int
	}{
		{
			name:      "small grid (5x5) wrap right",
			gridSize:  5,
			headX:     4,
			headY:     2,
			direction: direction.Right,
			wantX:     0,
			wantY:     2,
		},
		{
			name:      "large grid (20x20) wrap left",
			gridSize:  20,
			headX:     0,
			headY:     10,
			direction: direction.Left,
			wantX:     19,
			wantY:     10,
		},
		{
			name:      "single cell grid wrap",
			gridSize:  1,
			headX:     0,
			headY:     0,
			direction: direction.Right,
			wantX:     0,
			wantY:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.GameConfig{
				GridSize:    tt.gridSize,
				CellDimPx:   16,
				TargetScore: 100,
			}
			engine := NewPhysicsEngine(cfg)

			snake := &entities.Snake{
				Links: []*entities.Link{
					{X: tt.headX, Y: tt.headY},
				},
				Direction: tt.direction,
			}

			gotX, gotY := engine.CalculateNewHeadPosition(snake, tt.direction)

			if gotX != tt.wantX || gotY != tt.wantY {
				t.Errorf("CalculateNewHeadPosition() = (%d, %d), want (%d, %d)",
					gotX, gotY, tt.wantX, tt.wantY)
			}
		})
	}
}
