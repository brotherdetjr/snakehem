package gamestate

import (
	"snakehem/direction"
	"snakehem/internal/config"
	"snakehem/internal/entities"
	"snakehem/internal/interfaces"
)

// GameState represents a state in the game state machine
type GameState interface {
	Update(ctx *StateContext) (GameState, error)
	Name() string
}

// StateContext holds all the context needed for state updates
type StateContext struct {
	// Game entities
	Snakes []*entities.Snake
	Grid   *entities.GameGrid

	// Controllers (map of controllerID -> ControllerInput)
	Controllers map[string]interfaces.ControllerInput

	// Engines
	Physics PhysicsEngine
	Scoring ScoringEngine

	// Dependencies
	Random interfaces.RandomSource
	Config *config.GameConfig

	// State tracking
	Countdown     *int
	ElapsedFrames *uint64
	FadeCountdown *int
	ApplePresent  *bool
}

// PhysicsEngine handles movement and collision calculations
type PhysicsEngine interface {
	CalculateNewHeadPosition(snake *entities.Snake, dir direction.Direction) (x, y int)
}

// ScoringEngine handles score calculations
type ScoringEngine interface {
	ProcessBite(biter, victim *entities.Snake, linkIndex int, grid *entities.GameGrid, fadeCountdown *int) int
	ProcessApple(snake *entities.Snake, fadeCountdown *int) int
}
