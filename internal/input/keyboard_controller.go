package input

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"snakehem/internal/interfaces"
)

// KeyMapping defines the key bindings for a keyboard controller
type KeyMapping struct {
	Up    []ebiten.Key
	Down  []ebiten.Key
	Left  []ebiten.Key
	Right []ebiten.Key
	Start []ebiten.Key
	Exit  []ebiten.Key
}

// ArrowKeyMapping is the default arrow keys + numpad mapping
var ArrowKeyMapping = KeyMapping{
	Up:    []ebiten.Key{ebiten.KeyArrowUp, ebiten.KeyNumpad8},
	Down:  []ebiten.Key{ebiten.KeyArrowDown, ebiten.KeyNumpad2},
	Left:  []ebiten.Key{ebiten.KeyArrowLeft, ebiten.KeyNumpad4},
	Right: []ebiten.Key{ebiten.KeyArrowRight, ebiten.KeyNumpad6},
	Start: []ebiten.Key{ebiten.KeySpace, ebiten.KeyEnter, ebiten.KeyNumpadEnter},
	Exit:  []ebiten.Key{ebiten.KeyEscape},
}

// WASDKeyMapping is the WASD key mapping
var WASDKeyMapping = KeyMapping{
	Up:    []ebiten.Key{ebiten.KeyW},
	Down:  []ebiten.Key{ebiten.KeyS},
	Left:  []ebiten.Key{ebiten.KeyA},
	Right: []ebiten.Key{ebiten.KeyD},
	Start: []ebiten.Key{ebiten.KeyTab},
	Exit:  []ebiten.Key{ebiten.KeyGraveAccent},
}

// KeyboardController implements ControllerInput for keyboard with configurable key mappings
type KeyboardController struct {
	id      string
	mapping KeyMapping
}

// NewKeyboardController creates a new keyboard controller with the given ID and key mapping
func NewKeyboardController(id string, mapping KeyMapping) *KeyboardController {
	return &KeyboardController{
		id:      id,
		mapping: mapping,
	}
}

// ID returns the controller identifier
func (k *KeyboardController) ID() string {
	return k.id
}

// Equals checks if this controller is the same as another controller
func (k *KeyboardController) Equals(other interfaces.ControllerInput) bool {
	otherKb, ok := other.(*KeyboardController)
	if !ok {
		return false
	}
	return k.id == otherKb.id
}

// IsAnyJustPressed checks if any of this controller's keys were just pressed
func (k *KeyboardController) IsAnyJustPressed() bool {
	return k.IsUpJustPressed() || k.IsDownJustPressed() || k.IsLeftJustPressed() ||
		k.IsRightJustPressed() || k.IsExitJustPressed() || k.IsStartJustPressed()
}

// IsUpJustPressed checks if the up key was just pressed
func (k *KeyboardController) IsUpJustPressed() bool {
	return k.anyKeyJustPressed(k.mapping.Up)
}

// IsDownJustPressed checks if the down key was just pressed
func (k *KeyboardController) IsDownJustPressed() bool {
	return k.anyKeyJustPressed(k.mapping.Down)
}

// IsLeftJustPressed checks if the left key was just pressed
func (k *KeyboardController) IsLeftJustPressed() bool {
	return k.anyKeyJustPressed(k.mapping.Left)
}

// IsRightJustPressed checks if the right key was just pressed
func (k *KeyboardController) IsRightJustPressed() bool {
	return k.anyKeyJustPressed(k.mapping.Right)
}

// IsStartJustPressed checks if the start key was just pressed
func (k *KeyboardController) IsStartJustPressed() bool {
	return k.anyKeyJustPressed(k.mapping.Start)
}

// IsExitJustPressed checks if the exit key was just pressed
func (k *KeyboardController) IsExitJustPressed() bool {
	return k.anyKeyJustPressed(k.mapping.Exit)
}

// Vibrate does nothing for keyboard (no haptic feedback)
func (k *KeyboardController) Vibrate(_ time.Duration) {
	// Keyboards don't vibrate
}

// anyKeyJustPressed checks if any key in the list was just pressed
func (k *KeyboardController) anyKeyJustPressed(keys []ebiten.Key) bool {
	for _, key := range keys {
		if inpututil.IsKeyJustPressed(key) {
			return true
		}
	}
	return false
}

// Ensure KeyboardController implements interfaces.ControllerInput
var _ interfaces.ControllerInput = (*KeyboardController)(nil)
