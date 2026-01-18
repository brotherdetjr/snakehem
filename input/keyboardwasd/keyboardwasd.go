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

func (k keyboardWasd) IsUpJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyW)
}

func (k keyboardWasd) IsDownJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyS)
}

func (k keyboardWasd) IsLeftJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyA)
}

func (k keyboardWasd) IsRightJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyD)
}

func (k keyboardWasd) IsExitJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyGraveAccent)
}

func (k keyboardWasd) IsStartJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyTab)
}

func (k keyboardWasd) Vibrate(_ time.Duration) {
	// nothing
}
