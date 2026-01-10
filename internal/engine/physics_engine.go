package engine

import (
	"snakehem/direction"
	"snakehem/internal/config"
	"snakehem/internal/entities"
)

// PhysicsEngine handles movement and collision calculations
type PhysicsEngine struct {
	config *config.GameConfig
}

// NewPhysicsEngine creates a new physics engine
func NewPhysicsEngine(config *config.GameConfig) *PhysicsEngine {
	return &PhysicsEngine{config: config}
}

// CalculateNewHeadPosition calculates where the snake's head will be after moving in the given direction
// Handles grid wrapping (toroidal grid)
func (p *PhysicsEngine) CalculateNewHeadPosition(snake *entities.Snake, dir direction.Direction) (x, y int) {
	head := snake.Links[0]
	nX := head.X + dir.Dx()
	nY := head.Y + dir.Dy()

	// Wrap around grid edges (toroidal grid)
	// Assuming Dx and Dy can only be -1, 0, 1
	if nX < 0 {
		nX = p.config.GridSize - 1
	}
	if nY < 0 {
		nY = p.config.GridSize - 1
	}
	if nX >= p.config.GridSize {
		nX = 0
	}
	if nY >= p.config.GridSize {
		nY = 0
	}

	return nX, nY
}
