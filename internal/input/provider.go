package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"snakehem/internal/interfaces"
)

// EbitenInputProvider manages all input controllers (keyboards and gamepads)
type EbitenInputProvider struct {
	keyboardArrows *KeyboardController
	keyboardWASD   *KeyboardController
	gamepads       map[ebiten.GamepadID]*GamepadController
}

// NewEbitenInputProvider creates a new input provider with default keyboard controllers
func NewEbitenInputProvider() *EbitenInputProvider {
	return &EbitenInputProvider{
		keyboardArrows: NewKeyboardController("keyboard-arrows", ArrowKeyMapping),
		keyboardWASD:   NewKeyboardController("keyboard-wasd", WASDKeyMapping),
		gamepads:       make(map[ebiten.GamepadID]*GamepadController),
	}
}

// GetControllerInputs returns all currently active controllers
// This includes keyboards and any connected gamepads
func (p *EbitenInputProvider) GetControllerInputs() []interfaces.ControllerInput {
	// Start with keyboard controllers
	controllers := []interfaces.ControllerInput{
		p.keyboardArrows,
		p.keyboardWASD,
	}

	// Update gamepad list (detect new/removed gamepads)
	p.updateGamepads()

	// Add all gamepad controllers
	for _, gamepad := range p.gamepads {
		controllers = append(controllers, gamepad)
	}

	return controllers
}

// IsGlobalExitPressed checks if the global exit key (Escape) was pressed
func (p *EbitenInputProvider) IsGlobalExitPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEscape)
}

// updateGamepads checks for new or disconnected gamepads
func (p *EbitenInputProvider) updateGamepads() {
	// Get list of currently connected gamepads
	connectedIDs := ebiten.AppendGamepadIDs(nil)

	// Create a set of connected IDs for quick lookup
	connected := make(map[ebiten.GamepadID]bool)
	for _, id := range connectedIDs {
		connected[id] = true

		// Add new gamepads
		if _, exists := p.gamepads[id]; !exists {
			p.gamepads[id] = NewGamepadController(id)
		}
	}

	// Remove disconnected gamepads
	for id := range p.gamepads {
		if !connected[id] {
			delete(p.gamepads, id)
		}
	}
}

// Ensure EbitenInputProvider implements interfaces.InputProvider
var _ interfaces.InputProvider = (*EbitenInputProvider)(nil)
