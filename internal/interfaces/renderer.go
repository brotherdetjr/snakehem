package interfaces

// Renderer handles all drawing operations
// This separates rendering from game logic
type Renderer interface {
	DrawBackground(screen Screen)
	DrawGrid(screen Screen, grid Grid)
	DrawUI(screen Screen, state GameState, snakeCount int, countdown int, elapsedFrames uint64)
	ApplyPostProcessing(screen Screen)
}

// Grid provides read-only access to the game grid
type Grid interface {
	Get(x, y int) interface{}
	Size() int
}

// GameState represents the current game state
type GameState interface {
	Name() string
}
