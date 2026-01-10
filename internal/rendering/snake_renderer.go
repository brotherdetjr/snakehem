package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"

	"snakehem/direction"
	"snakehem/internal/config"
	"snakehem/internal/entities"
)

// SnakeRenderer handles drawing snakes on screen
type SnakeRenderer struct {
	config *config.GameConfig
}

// NewSnakeRenderer creates a new snake renderer
func NewSnakeRenderer(config *config.GameConfig) *SnakeRenderer {
	return &SnakeRenderer{config: config}
}

// DrawSnake draws a single snake with all its links
func (sr *SnakeRenderer) DrawSnake(screen *ebiten.Image, snake *entities.Snake, countdown int) {
	for i, link := range snake.Links {
		isHead := i == 0
		sr.DrawLink(screen, link, snake, countdown, isHead)
	}
}

// DrawLink draws a single snake link
func (sr *SnakeRenderer) DrawLink(screen *ebiten.Image, link *entities.Link, snake *entities.Snake, countdown int, isHead bool) {
	shrink := (1 - float32(link.HealthPercent)/100) * float32(sr.config.CellDimPx) * 0.5

	if !isHead {
		// Draw regular body link
		vector.DrawFilledRect(
			screen,
			float32(link.X*sr.config.CellDimPx)+shrink,
			float32(link.Y*sr.config.CellDimPx)+shrink,
			float32(sr.config.CellDimPx)-shrink*2,
			float32(sr.config.CellDimPx)-shrink*2,
			withRedness(snake.Colour, link.Redness),
			false,
		)
		return
	}

	// Draw head
	if countdown > sr.config.Tps() {
		// During countdown, draw head as regular square
		vector.DrawFilledRect(
			screen,
			float32(link.X*sr.config.CellDimPx)+shrink,
			float32(link.Y*sr.config.CellDimPx)+shrink,
			float32(sr.config.CellDimPx)-shrink*2,
			float32(sr.config.CellDimPx)-shrink*2,
			withRedness(snake.Colour, link.Redness),
			false,
		)
	} else {
		// After countdown, draw head with directional eyes
		sr.DrawHeadWithEyes(screen, link, snake)
	}
}

// DrawHeadWithEyes draws the snake head with eyes pointing in the direction
func (sr *SnakeRenderer) DrawHeadWithEyes(screen *ebiten.Image, head *entities.Link, snake *entities.Snake) {
	x1, y1, x2, y2 := sr.calculateEyePositions(head, snake.Direction)

	if x1 != 0 || y1 != 0 || x2 != 0 || y2 != 0 {
		// Draw eyes
		vector.DrawFilledCircle(
			screen,
			x1,
			y1,
			float32(sr.config.EyeRadiusPx),
			withRedness(snake.Colour, head.Redness),
			false,
		)
		vector.DrawFilledCircle(
			screen,
			x2,
			y2,
			float32(sr.config.EyeRadiusPx),
			withRedness(snake.Colour, head.Redness),
			false,
		)
	} else {
		// No direction (direction.None), draw as regular square
		vector.DrawFilledRect(
			screen,
			float32(head.X*sr.config.CellDimPx),
			float32(head.Y*sr.config.CellDimPx),
			float32(sr.config.CellDimPx),
			float32(sr.config.CellDimPx),
			withRedness(snake.Colour, head.Redness),
			false,
		)
	}
}

// calculateEyePositions calculates eye positions based on direction
func (sr *SnakeRenderer) calculateEyePositions(head *entities.Link, dir direction.Direction) (x1, y1, x2, y2 float32) {
	eyeGap := float32(sr.config.EyeGapPx)
	cellDim := float32(sr.config.CellDimPx)

	switch dir {
	case direction.Up:
		x1 = float32(head.X)*cellDim + eyeGap
		y1 = float32(head.Y)*cellDim + eyeGap
		x2 = float32(head.X+1)*cellDim - eyeGap
		y2 = float32(head.Y)*cellDim + eyeGap
	case direction.Down:
		x1 = float32(head.X)*cellDim + eyeGap
		y1 = float32(head.Y+1)*cellDim - eyeGap
		x2 = float32(head.X+1)*cellDim - eyeGap
		y2 = float32(head.Y+1)*cellDim - eyeGap
	case direction.Left:
		x1 = float32(head.X)*cellDim + eyeGap
		y1 = float32(head.Y+1)*cellDim - eyeGap
		x2 = float32(head.X)*cellDim + eyeGap
		y2 = float32(head.Y)*cellDim + eyeGap
	case direction.Right:
		x1 = float32(head.X+1)*cellDim - eyeGap
		y1 = float32(head.Y+1)*cellDim - eyeGap
		x2 = float32(head.X+1)*cellDim - eyeGap
		y2 = float32(head.Y)*cellDim + eyeGap
	case direction.None:
		// Return zeros (caller will draw regular square)
		return 0, 0, 0, 0
	}

	return x1, y1, x2, y2
}

// DrawApple draws an apple on the grid
func (sr *SnakeRenderer) DrawApple(screen *ebiten.Image, apple *entities.Apple) {
	vector.DrawFilledRect(
		screen,
		float32(apple.X*sr.config.CellDimPx),
		float32(apple.Y*sr.config.CellDimPx),
		float32(sr.config.CellDimPx),
		float32(sr.config.CellDimPx),
		colornames.Red,
		false,
	)
}

// withRedness applies redness effect to a color (extracted helper function)
func withRedness(colour color.Color, redness float32) color.Color {
	red, green, blue, _ := colour.RGBA()
	r := float32(red >> 8)
	g := float32(green >> 8)
	b := float32(blue >> 8)
	return color.NRGBA{
		R: uint8(r + (255-r)*redness),
		G: uint8(g * (1 - redness)),
		B: uint8(b * (1 - redness)),
		A: 255,
	}
}
