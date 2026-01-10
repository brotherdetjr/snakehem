package interfaces

import (
	"image"
	"image/color"
	"time"
)

// GameEngine abstracts the game engine's RunGame method
// This allows us to mock the Ebiten engine for testing
type GameEngine interface {
	RunGame(game Game) error
}

// Game represents the core game interface
// This is similar to ebiten.Game but decoupled from Ebiten
type Game interface {
	Update() error
	Draw(screen Screen)
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
}

// Screen abstracts the drawing surface (ebiten.Image)
// This allows us to mock rendering for testing
type Screen interface {
	Fill(clr color.Color)
	DrawImage(img Image, opts *DrawImageOptions)
	Bounds() image.Rectangle
	SubImage(r image.Rectangle) Screen
}

// Image represents a drawable image
type Image interface {
	Bounds() image.Rectangle
	At(x, y int) color.Color
}

// DrawImageOptions represents options for drawing an image
type DrawImageOptions struct {
	GeoM   GeoM
	ColorM ColorM
}

// GeoM represents geometric transformation matrix
type GeoM interface {
	Reset()
	Translate(tx, ty float64)
	Scale(x, y float64)
	Rotate(theta float64)
}

// ColorM represents color transformation matrix
type ColorM interface {
	Reset()
	Scale(r, g, b, a float64)
}

// InputProvider provides access to all controller inputs
// This replaces the singleton controller pattern
type InputProvider interface {
	GetControllerInputs() []ControllerInput
	IsGlobalExitPressed() bool
}

// ControllerInput represents input from a single controller
// This replaces the controller.Controller interface
type ControllerInput interface {
	ID() string
	IsUpJustPressed() bool
	IsDownJustPressed() bool
	IsLeftJustPressed() bool
	IsRightJustPressed() bool
	IsStartJustPressed() bool
	IsExitJustPressed() bool
	IsAnyJustPressed() bool
	Vibrate(duration time.Duration)
	Equals(other ControllerInput) bool
}
