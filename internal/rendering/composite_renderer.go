package rendering

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"

	"snakehem/internal/config"
	"snakehem/internal/entities"
)

// CompositeRenderer composes all renderer components
// Note: Currently uses *ebiten.Image directly. Will be abstracted in Phase 7 (Ebiten adapter)
type CompositeRenderer struct {
	config        *config.GameConfig
	fontManager   *FontManager
	snakeRenderer *SnakeRenderer
	uiRenderer    *UIRenderer
	postProcessor *PostProcessor
}

// NewCompositeRenderer creates a new composite renderer with all components
func NewCompositeRenderer(
	config *config.GameConfig,
	fontManager *FontManager,
	snakeRenderer *SnakeRenderer,
	uiRenderer *UIRenderer,
	postProcessor *PostProcessor,
) *CompositeRenderer {
	return &CompositeRenderer{
		config:        config,
		fontManager:   fontManager,
		snakeRenderer: snakeRenderer,
		uiRenderer:    uiRenderer,
		postProcessor: postProcessor,
	}
}

// DrawBackground fills the screen with the background color
func (cr *CompositeRenderer) DrawBackground(screen *ebiten.Image) {
	screen.Fill(colornames.Darkolivegreen)
}

// DrawGrid draws all items on the grid (snakes and apples)
func (cr *CompositeRenderer) DrawGrid(screen *ebiten.Image, grid *entities.GameGrid, countdown int) {
	for i := 0; i < cr.config.GridSize; i++ {
		for j := 0; j < cr.config.GridSize; j++ {
			item := grid.Get(i, j)
			if item == nil {
				continue
			}

			switch v := item.(type) {
			case *entities.Link:
				snake := v.Snake
				isHead := snake.Links[0] == v
				cr.snakeRenderer.DrawLink(screen, v, snake, countdown, isHead)
			case *entities.Apple:
				cr.snakeRenderer.DrawApple(screen, v)
			}
		}
	}
}

// DrawLobbyUI draws the lobby state UI
func (cr *CompositeRenderer) DrawLobbyUI(screen *ebiten.Image, snakeCount int) {
	cr.uiRenderer.DrawLobbyUI(screen, snakeCount)
}

// DrawActionUI draws the action state UI (scores, countdown, time)
func (cr *CompositeRenderer) DrawActionUI(screen *ebiten.Image, snakes []*entities.Snake, countdown int, elapsedFrames uint64, fadeCountdown int) {
	// Draw fade overlay if someone reached target score
	if fadeCountdown > 0 {
		cr.uiRenderer.DrawFadeOverlay(screen, fadeCountdown)
	}

	// Draw scores
	cr.uiRenderer.DrawScores(screen, snakes, elapsedFrames, true, countdown)

	// Draw countdown
	cr.uiRenderer.DrawCountdown(screen, countdown, elapsedFrames)

	// Draw time elapsed
	cr.uiRenderer.DrawTimeElapsed(screen, elapsedFrames)
}

// DrawScoreboardUI draws the scoreboard state UI
func (cr *CompositeRenderer) DrawScoreboardUI(screen *ebiten.Image, snakes []*entities.Snake, elapsedFrames uint64) {
	cr.uiRenderer.DrawScoreboard(screen, snakes)
	cr.uiRenderer.DrawTimeElapsed(screen, elapsedFrames)
}

// ApplyPostProcessing applies the CRT shader effect
func (cr *CompositeRenderer) ApplyPostProcessing(screen *ebiten.Image) {
	cr.postProcessor.Apply(screen)
}
