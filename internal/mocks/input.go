package mocks

import (
	"time"

	"snakehem/internal/interfaces"
)

// MockControllerInput provides a mock controller for testing
type MockControllerInput struct {
	id               string
	upPressed        bool
	downPressed      bool
	leftPressed      bool
	rightPressed     bool
	startPressed     bool
	exitPressed      bool
	anyPressed       bool
	vibrateDuration  time.Duration
	vibrateCallCount int
}

// NewMockControllerInput creates a new mock controller with the given ID
func NewMockControllerInput(id string) *MockControllerInput {
	return &MockControllerInput{id: id}
}

// ID returns the controller ID
func (m *MockControllerInput) ID() string {
	return m.id
}

// IsUpJustPressed returns whether up was just pressed
func (m *MockControllerInput) IsUpJustPressed() bool {
	return m.upPressed
}

// IsDownJustPressed returns whether down was just pressed
func (m *MockControllerInput) IsDownJustPressed() bool {
	return m.downPressed
}

// IsLeftJustPressed returns whether left was just pressed
func (m *MockControllerInput) IsLeftJustPressed() bool {
	return m.leftPressed
}

// IsRightJustPressed returns whether right was just pressed
func (m *MockControllerInput) IsRightJustPressed() bool {
	return m.rightPressed
}

// IsStartJustPressed returns whether start was just pressed
func (m *MockControllerInput) IsStartJustPressed() bool {
	return m.startPressed
}

// IsExitJustPressed returns whether exit was just pressed
func (m *MockControllerInput) IsExitJustPressed() bool {
	return m.exitPressed
}

// IsAnyJustPressed returns whether any button was just pressed
func (m *MockControllerInput) IsAnyJustPressed() bool {
	return m.anyPressed
}

// Vibrate records the vibration duration
func (m *MockControllerInput) Vibrate(duration time.Duration) {
	m.vibrateDuration = duration
	m.vibrateCallCount++
}

// Equals checks if two controllers are equal
func (m *MockControllerInput) Equals(other interfaces.ControllerInput) bool {
	if other == nil {
		return false
	}
	return m.id == other.ID()
}

// SetUpPressed sets the up button state
func (m *MockControllerInput) SetUpPressed(pressed bool) {
	m.upPressed = pressed
}

// SetDownPressed sets the down button state
func (m *MockControllerInput) SetDownPressed(pressed bool) {
	m.downPressed = pressed
}

// SetLeftPressed sets the left button state
func (m *MockControllerInput) SetLeftPressed(pressed bool) {
	m.leftPressed = pressed
}

// SetRightPressed sets the right button state
func (m *MockControllerInput) SetRightPressed(pressed bool) {
	m.rightPressed = pressed
}

// SetStartPressed sets the start button state
func (m *MockControllerInput) SetStartPressed(pressed bool) {
	m.startPressed = pressed
}

// SetExitPressed sets the exit button state
func (m *MockControllerInput) SetExitPressed(pressed bool) {
	m.exitPressed = pressed
}

// SetAnyPressed sets the any button state
func (m *MockControllerInput) SetAnyPressed(pressed bool) {
	m.anyPressed = pressed
}

// GetVibrateCallCount returns how many times Vibrate was called
func (m *MockControllerInput) GetVibrateCallCount() int {
	return m.vibrateCallCount
}

// GetVibrateDuration returns the last vibrate duration
func (m *MockControllerInput) GetVibrateDuration() time.Duration {
	return m.vibrateDuration
}

// Verify interface compliance at compile time
var _ interfaces.ControllerInput = (*MockControllerInput)(nil)

// MockInputProvider provides mock input for testing
type MockInputProvider struct {
	controllers       []interfaces.ControllerInput
	globalExitPressed bool
}

// NewMockInputProvider creates a new mock input provider
func NewMockInputProvider(controllers ...interfaces.ControllerInput) *MockInputProvider {
	return &MockInputProvider{
		controllers: controllers,
	}
}

// GetControllerInputs returns the mock controllers
func (m *MockInputProvider) GetControllerInputs() []interfaces.ControllerInput {
	return m.controllers
}

// IsGlobalExitPressed returns whether global exit was pressed
func (m *MockInputProvider) IsGlobalExitPressed() bool {
	return m.globalExitPressed
}

// SetGlobalExitPressed sets the global exit state
func (m *MockInputProvider) SetGlobalExitPressed(pressed bool) {
	m.globalExitPressed = pressed
}

// AddController adds a controller to the provider
func (m *MockInputProvider) AddController(ctrl interfaces.ControllerInput) {
	m.controllers = append(m.controllers, ctrl)
}

// RemoveController removes a controller from the provider
func (m *MockInputProvider) RemoveController(id string) {
	for i, ctrl := range m.controllers {
		if ctrl.ID() == id {
			m.controllers = append(m.controllers[:i], m.controllers[i+1:]...)
			return
		}
	}
}

// ClearControllers removes all controllers
func (m *MockInputProvider) ClearControllers() {
	m.controllers = nil
}

// Verify interface compliance at compile time
var _ interfaces.InputProvider = (*MockInputProvider)(nil)
