package ebiten_adapter

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"snakehem/internal/interfaces"
)

// EbitenEngine wraps ebiten.RunGame to implement interfaces.GameEngine
type EbitenEngine struct{}

// NewEbitenEngine creates a new Ebiten engine adapter
func NewEbitenEngine() *EbitenEngine {
	return &EbitenEngine{}
}

// RunGame wraps ebiten.RunGame and adapts our Game interface
func (e *EbitenEngine) RunGame(game interfaces.Game) error {
	adapter := &gameAdapter{game: game}
	return ebiten.RunGame(adapter)
}

// gameAdapter adapts our interfaces.Game to ebiten.Game
type gameAdapter struct {
	game interfaces.Game
}

// Update implements ebiten.Game
func (a *gameAdapter) Update() error {
	return a.game.Update()
}

// Draw implements ebiten.Game
func (a *gameAdapter) Draw(screen *ebiten.Image) {
	screenAdapter := &ScreenAdapter{image: screen}
	a.game.Draw(screenAdapter)
}

// Layout implements ebiten.Game
func (a *gameAdapter) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return a.game.Layout(outsideWidth, outsideHeight)
}

// ScreenAdapter wraps *ebiten.Image to implement interfaces.Screen
type ScreenAdapter struct {
	image *ebiten.Image
}

// NewScreenAdapter creates a new screen adapter from an ebiten.Image
func NewScreenAdapter(img *ebiten.Image) *ScreenAdapter {
	return &ScreenAdapter{image: img}
}

// EbitenImage returns the underlying *ebiten.Image
// This is needed for renderers that use Ebiten-specific functions (like vector drawing)
func (s *ScreenAdapter) EbitenImage() *ebiten.Image {
	return s.image
}

// Fill implements interfaces.Screen
func (s *ScreenAdapter) Fill(clr color.Color) {
	s.image.Fill(clr)
}

// DrawImage implements interfaces.Screen
func (s *ScreenAdapter) DrawImage(img interfaces.Image, opts *interfaces.DrawImageOptions) {
	// Convert interfaces.Image to *ebiten.Image
	ebitenImg, ok := img.(*ebiten.Image)
	if !ok {
		// If not already an ebiten.Image, we can't draw it
		// This shouldn't happen in production code
		return
	}

	// Convert interfaces.DrawImageOptions to *ebiten.DrawImageOptions
	var ebitenOpts *ebiten.DrawImageOptions
	if opts != nil {
		ebitenOpts = &ebiten.DrawImageOptions{}

		// Convert GeoM
		if opts.GeoM != nil {
			geomAdapter, ok := opts.GeoM.(*GeoMAdapter)
			if ok {
				ebitenOpts.GeoM = geomAdapter.matrix
			}
		}

		// Convert ColorM
		if opts.ColorM != nil {
			colormAdapter, ok := opts.ColorM.(*ColorMAdapter)
			if ok {
				ebitenOpts.ColorM = colormAdapter.matrix
			}
		}
	}

	s.image.DrawImage(ebitenImg, ebitenOpts)
}

// Bounds implements interfaces.Screen
func (s *ScreenAdapter) Bounds() image.Rectangle {
	return s.image.Bounds()
}

// SubImage implements interfaces.Screen
func (s *ScreenAdapter) SubImage(r image.Rectangle) interfaces.Screen {
	subImg := s.image.SubImage(r).(*ebiten.Image)
	return &ScreenAdapter{image: subImg}
}

// GeoMAdapter wraps ebiten.GeoM to implement interfaces.GeoM
type GeoMAdapter struct {
	matrix ebiten.GeoM
}

// NewGeoM creates a new GeoM adapter
func NewGeoM() *GeoMAdapter {
	return &GeoMAdapter{}
}

// Reset implements interfaces.GeoM
func (g *GeoMAdapter) Reset() {
	g.matrix.Reset()
}

// Translate implements interfaces.GeoM
func (g *GeoMAdapter) Translate(tx, ty float64) {
	g.matrix.Translate(tx, ty)
}

// Scale implements interfaces.GeoM
func (g *GeoMAdapter) Scale(x, y float64) {
	g.matrix.Scale(x, y)
}

// Rotate implements interfaces.GeoM
func (g *GeoMAdapter) Rotate(theta float64) {
	g.matrix.Rotate(theta)
}

// ColorMAdapter wraps ebiten.ColorM to implement interfaces.ColorM
type ColorMAdapter struct {
	matrix ebiten.ColorM
}

// NewColorM creates a new ColorM adapter
func NewColorM() *ColorMAdapter {
	return &ColorMAdapter{}
}

// Reset implements interfaces.ColorM
func (c *ColorMAdapter) Reset() {
	c.matrix.Reset()
}

// Scale implements interfaces.ColorM
func (c *ColorMAdapter) Scale(r, g, b, a float64) {
	c.matrix.Scale(r, g, b, a)
}
