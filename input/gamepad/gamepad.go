package gamepad

import (
	"snakehem/input/controller"
	"snakehem/model"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

func (g Gamepad) IsAnyPressed() bool {
	return g.IsUpPressed() || g.IsDownPressed() || g.IsLeftPressed() ||
		g.IsRightPressed() || g.IsStartPressed() || g.IsExitPressed()
}

func (g Gamepad) IsUpJustPressed() bool {
	id := ebiten.GamepadID(g)
	return inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightTop) ||
		inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonLeftTop)
}

func (g Gamepad) IsUpPressed() bool {
	id := ebiten.GamepadID(g)
	durRightTop := inpututil.StandardGamepadButtonPressDuration(id, ebiten.StandardGamepadButtonRightTop)
	durLeftTop := inpututil.StandardGamepadButtonPressDuration(id, ebiten.StandardGamepadButtonLeftTop)
	return durRightTop > 0 && durRightTop%model.TpsMultiplier == 0 ||
		durLeftTop > 0 && durLeftTop%model.TpsMultiplier == 0
}

func (g Gamepad) IsDownJustPressed() bool {
	id := ebiten.GamepadID(g)
	return inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightBottom) ||
		inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonLeftBottom)
}

func (g Gamepad) IsDownPressed() bool {
	id := ebiten.GamepadID(g)
	durRightBottom := inpututil.StandardGamepadButtonPressDuration(id, ebiten.StandardGamepadButtonRightBottom)
	durLeftBottom := inpututil.StandardGamepadButtonPressDuration(id, ebiten.StandardGamepadButtonLeftBottom)
	return durRightBottom > 0 && durRightBottom%model.TpsMultiplier == 0 ||
		durLeftBottom > 0 && durLeftBottom%model.TpsMultiplier == 0
}

func (g Gamepad) IsLeftJustPressed() bool {
	id := ebiten.GamepadID(g)
	return inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightLeft) ||
		inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonLeftLeft)
}

func (g Gamepad) IsLeftPressed() bool {
	id := ebiten.GamepadID(g)
	durRightLeft := inpututil.StandardGamepadButtonPressDuration(id, ebiten.StandardGamepadButtonRightLeft)
	durLeftLeft := inpututil.StandardGamepadButtonPressDuration(id, ebiten.StandardGamepadButtonLeftLeft)
	return durRightLeft > 0 && durRightLeft%model.TpsMultiplier == 0 ||
		durLeftLeft > 0 && durLeftLeft%model.TpsMultiplier == 0
}

func (g Gamepad) IsRightJustPressed() bool {
	id := ebiten.GamepadID(g)
	return inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightRight) ||
		inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonLeftRight)
}

func (g Gamepad) IsRightPressed() bool {
	id := ebiten.GamepadID(g)
	durRightRight := inpututil.StandardGamepadButtonPressDuration(id, ebiten.StandardGamepadButtonRightRight)
	durLeftRight := inpututil.StandardGamepadButtonPressDuration(id, ebiten.StandardGamepadButtonLeftRight)
	return durRightRight > 0 && durRightRight%model.TpsMultiplier == 0 ||
		durLeftRight > 0 && durLeftRight%model.TpsMultiplier == 0
}

func (g Gamepad) IsExitJustPressed() bool {
	return inpututil.IsStandardGamepadButtonJustPressed(ebiten.GamepadID(g), ebiten.StandardGamepadButtonCenterLeft)
}

func (g Gamepad) IsExitPressed() bool {
	id := ebiten.GamepadID(g)
	dur := inpututil.StandardGamepadButtonPressDuration(id, ebiten.StandardGamepadButtonCenterLeft)
	return dur > 0 && dur%model.TpsMultiplier == 0
}

func (g Gamepad) IsStartJustPressed() bool {
	return inpututil.IsStandardGamepadButtonJustPressed(ebiten.GamepadID(g), ebiten.StandardGamepadButtonCenterRight)
}

func (g Gamepad) IsStartPressed() bool {
	id := ebiten.GamepadID(g)
	dur := inpututil.StandardGamepadButtonPressDuration(id, ebiten.StandardGamepadButtonCenterRight)
	return dur > 0 && dur%model.TpsMultiplier == 0
}

func (g Gamepad) Vibrate(duration time.Duration) {
	ebiten.VibrateGamepad(ebiten.GamepadID(g), &ebiten.VibrateGamepadOptions{
		Duration:        duration,
		StrongMagnitude: 1,
		WeakMagnitude:   1,
	})
}
