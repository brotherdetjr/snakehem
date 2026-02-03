package keyboard

import (
	"snakehem/input/controller"
	"snakehem/model"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type keyboard struct {
}

var Instance controller.Controller = keyboard{}

func (k keyboard) Equals(controller controller.Controller) bool {
	_, ok := controller.(keyboard)
	return ok
}

func (k keyboard) IsAnyJustPressed() bool {
	return k.IsUpJustPressed() || k.IsDownJustPressed() || k.IsLeftJustPressed() ||
		k.IsRightJustPressed() || k.IsExitJustPressed() || k.IsStartJustPressed()
}

func (k keyboard) IsAnyPressed() bool {
	return k.IsUpPressed() || k.IsDownPressed() || k.IsLeftPressed() ||
		k.IsRightPressed() || k.IsStartPressed() || k.IsExitPressed()
}

func (k keyboard) IsUpJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad8)
}

func (k keyboard) IsUpPressed() bool {
	durArrowUp := inpututil.KeyPressDuration(ebiten.KeyArrowUp)
	durNumpad8 := inpututil.KeyPressDuration(ebiten.KeyNumpad8)
	return durArrowUp > 0 && durArrowUp%model.TpsMultiplier == 0 ||
		durNumpad8 > 0 && durNumpad8%model.TpsMultiplier == 0
}

func (k keyboard) IsDownJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad2)
}

func (k keyboard) IsDownPressed() bool {
	durArrowDown := inpututil.KeyPressDuration(ebiten.KeyArrowDown)
	durNumpad2 := inpututil.KeyPressDuration(ebiten.KeyNumpad2)
	return durArrowDown > 0 && durArrowDown%model.TpsMultiplier == 0 ||
		durNumpad2 > 0 && durNumpad2%model.TpsMultiplier == 0
}

func (k keyboard) IsLeftJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad4)
}

func (k keyboard) IsLeftPressed() bool {
	durArrowLeft := inpututil.KeyPressDuration(ebiten.KeyArrowLeft)
	durNumpad4 := inpututil.KeyPressDuration(ebiten.KeyNumpad4)
	return durArrowLeft > 0 && durArrowLeft%model.TpsMultiplier == 0 ||
		durNumpad4 > 0 && durNumpad4%model.TpsMultiplier == 0
}

func (k keyboard) IsRightJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad6)
}

func (k keyboard) IsRightPressed() bool {
	durArrowRight := inpututil.KeyPressDuration(ebiten.KeyArrowRight)
	durNumpad6 := inpututil.KeyPressDuration(ebiten.KeyNumpad6)
	return durArrowRight > 0 && durArrowRight%model.TpsMultiplier == 0 ||
		durNumpad6 > 0 && durNumpad6%model.TpsMultiplier == 0
}

func (k keyboard) IsExitJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEscape)
}

func (k keyboard) IsExitPressed() bool {
	durEscape := inpututil.KeyPressDuration(ebiten.KeyEscape)
	return durEscape > 0 && durEscape%model.TpsMultiplier == 0
}

func (k keyboard) IsStartJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter)
}

func (k keyboard) IsStartPressed() bool {
	durSpace := inpututil.KeyPressDuration(ebiten.KeySpace)
	durEnter := inpututil.KeyPressDuration(ebiten.KeyEnter)
	durNumpadEnter := inpututil.KeyPressDuration(ebiten.KeyNumpadEnter)
	return durSpace > 0 && durSpace%model.TpsMultiplier == 0 ||
		durEnter > 0 && durEnter%model.TpsMultiplier == 0 ||
		durNumpadEnter > 0 && durNumpadEnter%model.TpsMultiplier == 0
}

func (k keyboard) Vibrate(_ time.Duration) {
	// nothing
}
