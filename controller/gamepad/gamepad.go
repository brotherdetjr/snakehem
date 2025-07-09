package gamepad

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"snakehem/controller"
	"time"
)

type Gamepad ebiten.GamepadID

func NewGamepad(id ebiten.GamepadID) Gamepad {
	return Gamepad(id)
}

func (g Gamepad) Equals(controller controller.Controller) bool {
	other, ok := controller.(Gamepad)
	return ok && g == other
}

func (g Gamepad) IsAnyJustPressed() bool {
	buttonPressed := false
	for b := ebiten.StandardGamepadButton(0); b <= ebiten.StandardGamepadButtonMax; b++ {
		if inpututil.IsStandardGamepadButtonJustPressed(ebiten.GamepadID(g), b) {
			buttonPressed = true
		}
	}
	return buttonPressed
}

func (g Gamepad) IsUpJustPressed() bool {
	id := ebiten.GamepadID(g)
	return inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightTop) ||
		inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonLeftTop)
}

func (g Gamepad) IsDownJustPressed() bool {
	id := ebiten.GamepadID(g)
	return inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightBottom) ||
		inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonLeftBottom)
}

func (g Gamepad) IsLeftJustPressed() bool {
	id := ebiten.GamepadID(g)
	return inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightLeft) ||
		inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonLeftLeft)
}

func (g Gamepad) IsRightJustPressed() bool {
	id := ebiten.GamepadID(g)
	return inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightRight) ||
		inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonLeftRight)
}

func (g Gamepad) IsExitJustPressed() bool {
	return inpututil.IsStandardGamepadButtonJustPressed(ebiten.GamepadID(g), ebiten.StandardGamepadButtonCenterLeft)
}

func (g Gamepad) IsStartJustPressed() bool {
	return inpututil.IsStandardGamepadButtonJustPressed(ebiten.GamepadID(g), ebiten.StandardGamepadButtonCenterRight)
}

func (g Gamepad) Vibrate(duration time.Duration) {
	ebiten.VibrateGamepad(ebiten.GamepadID(g), &ebiten.VibrateGamepadOptions{
		Duration:        duration,
		StrongMagnitude: 1,
		WeakMagnitude:   1,
	})
}
