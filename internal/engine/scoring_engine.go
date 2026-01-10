package engine

import (
	"snakehem/internal/config"
	"snakehem/internal/entities"
)

// ScoringEngine handles all score calculations
type ScoringEngine struct {
	config *config.GameConfig
}

// NewScoringEngine creates a new scoring engine
func NewScoringEngine(config *config.GameConfig) *ScoringEngine {
	return &ScoringEngine{config: config}
}

// ProcessBite calculates points for biting another snake
// Returns the points earned and updates fadeCountdown if target score is reached
func (p *ScoringEngine) ProcessBite(biter, victim *entities.Snake, linkIndex int, grid *entities.GameGrid, fadeCountdown *int) int {
	points := 0

	// Points for biting a link
	if victim != biter {
		points += p.config.BitLinkScore
	}

	// Bonus points if the tail is nipped (link health depleted)
	bittenLink := victim.Links[linkIndex]
	if bittenLink.HealthPercent <= 0 {
		if victim != biter {
			// Bonus for nipping tail: points for each remaining link
			remainingLinks := len(victim.Links) - linkIndex + 1
			points += remainingLinks * p.config.NippedTailLinkBonusMultiplier
		}
	}

	return points
}

// ProcessApple calculates points for eating an apple
// Updates fadeCountdown if target score is reached
func (p *ScoringEngine) ProcessApple(snake *entities.Snake, fadeCountdown *int) int {
	points := p.config.AppleScore

	// Check if this puts the snake at or above target score
	if snake.Score+points >= p.config.TargetScore {
		*fadeCountdown = p.config.GridFadeCountdown
	}

	return points
}

// HasWinner checks if any snake has reached the target score
func (p *ScoringEngine) HasWinner(snakes []*entities.Snake) bool {
	for _, snake := range snakes {
		if snake.Score >= p.config.TargetScore {
			return true
		}
	}
	return false
}
