package engine

import (
	"testing"

	"snakehem/internal/config"
	"snakehem/internal/entities"
)

func TestScoringEngine_ProcessApple(t *testing.T) {
	cfg := &config.GameConfig{
		GridSize:          10,
		CellDimPx:         16,
		TargetScore:       100,
		AppleScore:        10,
		GridFadeCountdown: 30,
	}
	engine := NewScoringEngine(cfg)

	tests := []struct {
		name              string
		snakeScore        int
		wantPoints        int
		wantFadeCountdown int
	}{
		{
			name:              "normal apple consumption",
			snakeScore:        50,
			wantPoints:        10,
			wantFadeCountdown: 0,
		},
		{
			name:              "apple consumption reaches target",
			snakeScore:        90,
			wantPoints:        10,
			wantFadeCountdown: 30,
		},
		{
			name:              "apple consumption exceeds target",
			snakeScore:        95,
			wantPoints:        10,
			wantFadeCountdown: 30,
		},
		{
			name:              "first apple",
			snakeScore:        0,
			wantPoints:        10,
			wantFadeCountdown: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snake := &entities.Snake{
				Score: tt.snakeScore,
			}
			fadeCountdown := 0

			points := engine.ProcessApple(snake, &fadeCountdown)

			if points != tt.wantPoints {
				t.Errorf("ProcessApple() points = %d, want %d", points, tt.wantPoints)
			}

			if fadeCountdown != tt.wantFadeCountdown {
				t.Errorf("ProcessApple() fadeCountdown = %d, want %d", fadeCountdown, tt.wantFadeCountdown)
			}
		})
	}
}

func TestScoringEngine_ProcessBite(t *testing.T) {
	cfg := &config.GameConfig{
		GridSize:                      10,
		CellDimPx:                     16,
		TargetScore:                   100,
		BitLinkScore:                  5,
		NippedTailLinkBonusMultiplier: 2,
	}
	engine := NewScoringEngine(cfg)

	tests := []struct {
		name            string
		biterSame       bool // if true, biter == victim (self-bite)
		linkIndex       int
		linkHealth      int8
		victimLinkCount int
		wantPoints      int
	}{
		{
			name:            "bite other snake's healthy link",
			biterSame:       false,
			linkIndex:       2,
			linkHealth:      100,
			victimLinkCount: 5,
			wantPoints:      5, // just BitLinkScore
		},
		{
			name:            "bite other snake's depleted link (nip tail)",
			biterSame:       false,
			linkIndex:       2,
			linkHealth:      0,
			victimLinkCount: 5,
			wantPoints:      5 + (5-2+1)*2, // BitLinkScore + (remaining links * multiplier)
		},
		{
			name:            "bite self (no points)",
			biterSame:       true,
			linkIndex:       1,
			linkHealth:      100,
			victimLinkCount: 3,
			wantPoints:      0,
		},
		{
			name:            "bite self's depleted link (no bonus)",
			biterSame:       true,
			linkIndex:       1,
			linkHealth:      0,
			victimLinkCount: 3,
			wantPoints:      0,
		},
		{
			name:            "nip tail at last link",
			biterSame:       false,
			linkIndex:       4,
			linkHealth:      0,
			victimLinkCount: 5,
			wantPoints:      5 + (5-4+1)*2, // BitLinkScore + 2 links * 2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create victim snake
			victim := &entities.Snake{
				Links: make([]*entities.Link, tt.victimLinkCount),
			}
			for i := range victim.Links {
				victim.Links[i] = &entities.Link{
					X:             i,
					Y:             0,
					HealthPercent: 100,
				}
			}
			victim.Links[tt.linkIndex].HealthPercent = tt.linkHealth

			// Create biter snake (same as victim if self-bite)
			var biter *entities.Snake
			if tt.biterSame {
				biter = victim
			} else {
				biter = &entities.Snake{
					Links: []*entities.Link{{X: 0, Y: 1}},
				}
			}

			grid := entities.NewGameGrid(cfg.GridSize)
			fadeCountdown := 0

			points := engine.ProcessBite(biter, victim, tt.linkIndex, grid, &fadeCountdown)

			if points != tt.wantPoints {
				t.Errorf("ProcessBite() = %d, want %d", points, tt.wantPoints)
			}
		})
	}
}

func TestScoringEngine_HasWinner(t *testing.T) {
	cfg := &config.GameConfig{
		GridSize:    10,
		CellDimPx:   16,
		TargetScore: 100,
	}
	engine := NewScoringEngine(cfg)

	tests := []struct {
		name       string
		scores     []int
		wantWinner bool
	}{
		{
			name:       "no winner",
			scores:     []int{50, 60, 70},
			wantWinner: false,
		},
		{
			name:       "one winner exactly at target",
			scores:     []int{50, 100, 70},
			wantWinner: true,
		},
		{
			name:       "one winner above target",
			scores:     []int{50, 120, 70},
			wantWinner: true,
		},
		{
			name:       "multiple winners",
			scores:     []int{100, 110, 105},
			wantWinner: true,
		},
		{
			name:       "empty snake list",
			scores:     []int{},
			wantWinner: false,
		},
		{
			name:       "single snake no winner",
			scores:     []int{50},
			wantWinner: false,
		},
		{
			name:       "single snake winner",
			scores:     []int{100},
			wantWinner: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snakes := make([]*entities.Snake, len(tt.scores))
			for i, score := range tt.scores {
				snakes[i] = &entities.Snake{
					Score: score,
				}
			}

			got := engine.HasWinner(snakes)

			if got != tt.wantWinner {
				t.Errorf("HasWinner() = %v, want %v", got, tt.wantWinner)
			}
		})
	}
}
