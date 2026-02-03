package controller

import (
	"snakehem/model"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Controller interface {
	Equals(controller Controller) bool
	IsAnyJustPressed() bool
	IsAnyPressed() bool
	IsUpJustPressed() bool
	IsUpPressed() bool
	IsDownJustPressed() bool
	IsDownPressed() bool
	IsLeftJustPressed() bool
	IsLeftPressed() bool
	IsRightJustPressed() bool
	IsRightPressed() bool
	IsExitJustPressed() bool
	IsExitPressed() bool
	IsStartJustPressed() bool
	IsStartPressed() bool
	Vibrate(duration time.Duration)
}

func IsRepeatingGamepad(id ebiten.GamepadID, buttons ...ebiten.StandardGamepadButton) bool {
	for _, btn := range buttons {
		if inpututil.IsStandardGamepadButtonJustPressed(id, btn) {
			return true
		}
		dur := inpututil.StandardGamepadButtonPressDuration(id, btn)
		if dur > 0 && dur%model.ControllerRepeatPeriod == 0 {
			return true
		}
	}
	return false
}

func IsRepeatingKeyboard(keys ...ebiten.Key) bool {
	for _, key := range keys {
		if inpututil.IsKeyJustPressed(key) {
			return true
		}
		dur := inpututil.KeyPressDuration(key)
		if dur > 0 && dur%model.ControllerRepeatPeriod == 0 {
			return true
		}
	}
	return false
}
