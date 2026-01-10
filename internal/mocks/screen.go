package mocks

import (
	"image"
	"image/color"

	"snakehem/internal/interfaces"
)

// MockScreen provides a mock screen for testing rendering
type MockScreen struct {
	fillColor     color.Color
	fillCallCount int
	drawCalls     []DrawCall
	bounds        image.Rectangle
}

// DrawCall records a single DrawImage call
type DrawCall struct {
	Image interfaces.Image
	Opts  *interfaces.DrawImageOptions
}

// NewMockScreen creates a new mock screen with the given bounds
func NewMockScreen(width, height int) *MockScreen {
	return &MockScreen{
		bounds:    image.Rect(0, 0, width, height),
		drawCalls: make([]DrawCall, 0),
	}
}

// Fill records the fill color
func (m *MockScreen) Fill(clr color.Color) {
	m.fillColor = clr
	m.fillCallCount++
}

// DrawImage records the draw call
func (m *MockScreen) DrawImage(img interfaces.Image, opts *interfaces.DrawImageOptions) {
	m.drawCalls = append(m.drawCalls, DrawCall{
		Image: img,
		Opts:  opts,
	})
}

// Bounds returns the screen bounds
func (m *MockScreen) Bounds() image.Rectangle {
	return m.bounds
}

// SubImage creates a sub-image of the screen
func (m *MockScreen) SubImage(r image.Rectangle) interfaces.Screen {
	return &MockScreen{
		bounds:    r,
		drawCalls: make([]DrawCall, 0),
	}
}

// GetFillColor returns the last fill color used
func (m *MockScreen) GetFillColor() color.Color {
	return m.fillColor
}

// GetFillCallCount returns how many times Fill was called
func (m *MockScreen) GetFillCallCount() int {
	return m.fillCallCount
}

// GetDrawCalls returns all recorded draw calls
func (m *MockScreen) GetDrawCalls() []DrawCall {
	return m.drawCalls
}

// GetDrawCallCount returns how many times DrawImage was called
func (m *MockScreen) GetDrawCallCount() int {
	return len(m.drawCalls)
}

// Reset clears all recorded calls
func (m *MockScreen) Reset() {
	m.fillColor = nil
	m.fillCallCount = 0
	m.drawCalls = make([]DrawCall, 0)
}

// Verify interface compliance at compile time
var _ interfaces.Screen = (*MockScreen)(nil)

// MockGeoM provides a mock geometric transformation matrix
type MockGeoM struct {
	translateX, translateY float64
	scaleX, scaleY         float64
	rotateAngle            float64
	resetCalled            bool
}

// NewMockGeoM creates a new mock GeoM
func NewMockGeoM() *MockGeoM {
	return &MockGeoM{}
}

// Reset resets the transformation
func (m *MockGeoM) Reset() {
	m.resetCalled = true
	m.translateX = 0
	m.translateY = 0
	m.scaleX = 1
	m.scaleY = 1
	m.rotateAngle = 0
}

// Translate records the translation
func (m *MockGeoM) Translate(tx, ty float64) {
	m.translateX = tx
	m.translateY = ty
}

// Scale records the scale
func (m *MockGeoM) Scale(x, y float64) {
	m.scaleX = x
	m.scaleY = y
}

// Rotate records the rotation
func (m *MockGeoM) Rotate(theta float64) {
	m.rotateAngle = theta
}

// GetTranslate returns the translation values
func (m *MockGeoM) GetTranslate() (float64, float64) {
	return m.translateX, m.translateY
}

// GetScale returns the scale values
func (m *MockGeoM) GetScale() (float64, float64) {
	return m.scaleX, m.scaleY
}

// GetRotate returns the rotation angle
func (m *MockGeoM) GetRotate() float64 {
	return m.rotateAngle
}

// Verify interface compliance at compile time
var _ interfaces.GeoM = (*MockGeoM)(nil)

// MockColorM provides a mock color transformation matrix
type MockColorM struct {
	r, g, b, a  float64
	resetCalled bool
}

// NewMockColorM creates a new mock ColorM
func NewMockColorM() *MockColorM {
	return &MockColorM{
		r: 1, g: 1, b: 1, a: 1,
	}
}

// Reset resets the color transformation
func (m *MockColorM) Reset() {
	m.resetCalled = true
	m.r = 1
	m.g = 1
	m.b = 1
	m.a = 1
}

// Scale records the color scale
func (m *MockColorM) Scale(r, g, b, a float64) {
	m.r = r
	m.g = g
	m.b = b
	m.a = a
}

// GetScale returns the color scale values
func (m *MockColorM) GetScale() (float64, float64, float64, float64) {
	return m.r, m.g, m.b, m.a
}

// Verify interface compliance at compile time
var _ interfaces.ColorM = (*MockColorM)(nil)
