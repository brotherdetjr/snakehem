package keyboardwasd

import (
	"snakehem/input/controller"
	"snakehem/model"
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
	return k.IsUpPressed() || k.IsDownPressed() || k.IsLeftPressed() ||
		k.IsRightPressed() || k.IsStartPressed() || k.IsExitPressed()
}

func (k keyboardWasd) IsUpJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyW)
}

func (k keyboardWasd) IsUpPressed() bool {
	durW := inpututil.KeyPressDuration(ebiten.KeyW)
	return durW > 0 && durW%model.ControllerRepeatPeriod == 0
}

func (k keyboardWasd) IsDownJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyS)
}

func (k keyboardWasd) IsDownPressed() bool {
	durS := inpututil.KeyPressDuration(ebiten.KeyS)
	return durS > 0 && durS%model.ControllerRepeatPeriod == 0
}

func (k keyboardWasd) IsLeftJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyA)
}

func (k keyboardWasd) IsLeftPressed() bool {
	durA := inpututil.KeyPressDuration(ebiten.KeyA)
	return durA > 0 && durA%model.ControllerRepeatPeriod == 0
}

func (k keyboardWasd) IsRightJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyD)
}

func (k keyboardWasd) IsRightPressed() bool {
	durD := inpututil.KeyPressDuration(ebiten.KeyD)
	return durD > 0 && durD%model.ControllerRepeatPeriod == 0
}

func (k keyboardWasd) IsExitJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyGraveAccent)
}

func (k keyboardWasd) IsExitPressed() bool {
	durGraveAccent := inpututil.KeyPressDuration(ebiten.KeyGraveAccent)
	return durGraveAccent > 0 && durGraveAccent%model.ControllerRepeatPeriod == 0
}

func (k keyboardWasd) IsStartJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyTab)
}

func (k keyboardWasd) IsStartPressed() bool {
	durTab := inpututil.KeyPressDuration(ebiten.KeyTab)
	return durTab > 0 && durTab%model.ControllerRepeatPeriod == 0
}

func (k keyboardWasd) Vibrate(_ time.Duration) {
	// nothing
}
