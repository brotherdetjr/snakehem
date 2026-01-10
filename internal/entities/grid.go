package entities

import (
	"snakehem/internal/interfaces"
)

// GameGrid manages the spatial grid for game entities
type GameGrid struct {
	size  int
	cells [][]interface{}
}

// NewGameGrid creates a new game grid with the specified size
func NewGameGrid(size int) *GameGrid {
	cells := make([][]interface{}, size)
	for i := range cells {
		cells[i] = make([]interface{}, size)
	}
	return &GameGrid{
		size:  size,
		cells: cells,
	}
}

// Get returns the item at the specified coordinates
func (g *GameGrid) Get(x, y int) interface{} {
	if x < 0 || x >= g.size || y < 0 || y >= g.size {
		return nil
	}
	return g.cells[y][x]
}

// Set places an item at the specified coordinates
func (g *GameGrid) Set(x, y int, item interface{}) {
	if x >= 0 && x < g.size && y >= 0 && y < g.size {
		g.cells[y][x] = item
	}
}

// Clear removes the item at the specified coordinates
func (g *GameGrid) Clear(x, y int) {
	if x >= 0 && x < g.size && y >= 0 && y < g.size {
		g.cells[y][x] = nil
	}
}

// Size returns the grid size
func (g *GameGrid) Size() int {
	return g.size
}

// FindRandomEmpty finds a random empty cell in the grid
// Returns (x, y, true) if found, (-1, -1, false) if grid is full
func (g *GameGrid) FindRandomEmpty(random interfaces.RandomSource) (int, int, bool) {
	x := random.IntN(g.size)
	y := random.IntN(g.size)

	// Linear search from random starting point
	for yi := 0; yi < g.size; yi++ {
		for xi := 0; xi < g.size; xi++ {
			checkX := (x + xi) % g.size
			checkY := (y + yi) % g.size
			if g.cells[checkY][checkX] == nil {
				return checkX, checkY, true
			}
		}
	}

	return -1, -1, false
}

// Clear all cells in the grid
func (g *GameGrid) ClearAll() {
	for y := 0; y < g.size; y++ {
		for x := 0; x < g.size; x++ {
			g.cells[y][x] = nil
		}
	}
}
