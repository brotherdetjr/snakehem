package keyboardwasd

import (
	"snakehem/input/controller"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type keyboardWasd struct {
}

var Instance controller.Controller = keyboardWasd{}

func (k keyboardWasd) Equals(controller controller.Controller) bool {
	_, ok := controller.(keyboardWasd)
	return ok
}

func (k keyboardWasd) IsAnyJustPressed() bool {
	return k.IsUpJustPressed() || k.IsDownJustPressed() || k.IsLeftJustPressed() ||
		k.IsRightJustPressed() || k.IsExitJustPressed() || k.IsStartJustPressed()
}

func (k keyboardWasd) IsAnyPressed() bool {
	return k.IsAnyJustPressed() || k.IsUpPressed() || k.IsDownPressed() || k.IsLeftPressed() ||
		k.IsRightPressed() || k.IsStartPressed() || k.IsExitPressed()
}

func (k keyboardWasd) IsUpJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyW)
}

func (k keyboardWasd) IsUpPressed() bool {
	return controller.IsRepeatingKeyboard(ebiten.KeyW)
}

func (k keyboardWasd) IsDownJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyS)
}

func (k keyboardWasd) IsDownPressed() bool {
	return controller.IsRepeatingKeyboard(ebiten.KeyS)
}

func (k keyboardWasd) IsLeftJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyA)
}

func (k keyboardWasd) IsLeftPressed() bool {
	return controller.IsRepeatingKeyboard(ebiten.KeyA)
}

func (k keyboardWasd) IsRightJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyD)
}

func (k keyboardWasd) IsRightPressed() bool {
	return controller.IsRepeatingKeyboard(ebiten.KeyD)
}

func (k keyboardWasd) IsExitJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyGraveAccent)
}

func (k keyboardWasd) IsExitPressed() bool {
	return controller.IsRepeatingKeyboard(ebiten.KeyGraveAccent)
}

func (k keyboardWasd) IsStartJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyTab)
}

func (k keyboardWasd) IsStartPressed() bool {
	return controller.IsRepeatingKeyboard(ebiten.KeyTab)
}

func (k keyboardWasd) Vibrate(_ time.Duration) {
	// nothing
}
