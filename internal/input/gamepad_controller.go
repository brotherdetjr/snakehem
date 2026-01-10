package input

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"snakehem/internal/interfaces"
)

// GamepadController implements ControllerInput for a gamepad
type GamepadController struct {
	id       ebiten.GamepadID
	idString string
}

// NewGamepadController creates a new gamepad controller
func NewGamepadController(id ebiten.GamepadID) *GamepadController {
	return &GamepadController{
		id:       id,
		idString: fmt.Sprintf("gamepad-%d", id),
	}
}

// ID returns the controller identifier
func (g *GamepadController) ID() string {
	return g.idString
}

// Equals checks if this controller is the same as another controller
func (g *GamepadController) Equals(other interfaces.ControllerInput) bool {
	otherGp, ok := other.(*GamepadController)
	if !ok {
		return false
	}
	return g.id == otherGp.id
}

// IsAnyJustPressed checks if any button was just pressed
func (g *GamepadController) IsAnyJustPressed() bool {
	for b := ebiten.StandardGamepadButton(0); b <= ebiten.StandardGamepadButtonMax; b++ {
		if inpututil.IsStandardGamepadButtonJustPressed(g.id, b) {
			return true
		}
	}
	return false
}

// IsUpJustPressed checks if the up button/stick was just pressed
func (g *GamepadController) IsUpJustPressed() bool {
	return inpututil.IsStandardGamepadButtonJustPressed(g.id, ebiten.StandardGamepadButtonRightTop) ||
		inpututil.IsStandardGamepadButtonJustPressed(g.id, ebiten.StandardGamepadButtonLeftTop)
}

// IsDownJustPressed checks if the down button/stick was just pressed
func (g *GamepadController) IsDownJustPressed() bool {
	return inpututil.IsStandardGamepadButtonJustPressed(g.id, ebiten.StandardGamepadButtonRightBottom) ||
		inpututil.IsStandardGamepadButtonJustPressed(g.id, ebiten.StandardGamepadButtonLeftBottom)
}

// IsLeftJustPressed checks if the left button/stick was just pressed
func (g *GamepadController) IsLeftJustPressed() bool {
	return inpututil.IsStandardGamepadButtonJustPressed(g.id, ebiten.StandardGamepadButtonRightLeft) ||
		inpututil.IsStandardGamepadButtonJustPressed(g.id, ebiten.StandardGamepadButtonLeftLeft)
}

// IsRightJustPressed checks if the right button/stick was just pressed
func (g *GamepadController) IsRightJustPressed() bool {
	return inpututil.IsStandardGamepadButtonJustPressed(g.id, ebiten.StandardGamepadButtonRightRight) ||
		inpututil.IsStandardGamepadButtonJustPressed(g.id, ebiten.StandardGamepadButtonLeftRight)
}

// IsStartJustPressed checks if the start button was just pressed
func (g *GamepadController) IsStartJustPressed() bool {
	return inpututil.IsStandardGamepadButtonJustPressed(g.id, ebiten.StandardGamepadButtonCenterRight)
}

// IsExitJustPressed checks if the exit/select button was just pressed
func (g *GamepadController) IsExitJustPressed() bool {
	return inpututil.IsStandardGamepadButtonJustPressed(g.id, ebiten.StandardGamepadButtonCenterLeft)
}

// Vibrate activates gamepad vibration for the specified duration
func (g *GamepadController) Vibrate(duration time.Duration) {
	ebiten.VibrateGamepad(g.id, &ebiten.VibrateGamepadOptions{
		Duration:        duration,
		StrongMagnitude: 1,
		WeakMagnitude:   1,
	})
}

// Ensure GamepadController implements interfaces.ControllerInput
var _ interfaces.ControllerInput = (*GamepadController)(nil)
