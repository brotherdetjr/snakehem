package keyboard

import (
	"snakehem/input/controller"
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
	return k.IsAnyJustPressed() || k.IsUpPressed() || k.IsDownPressed() || k.IsLeftPressed() ||
		k.IsRightPressed() || k.IsStartPressed() || k.IsExitPressed()
}

func (k keyboard) IsUpJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad8)
}

func (k keyboard) IsUpPressed() bool {
	return controller.IsRepeatingKeyboard(ebiten.KeyArrowUp, ebiten.KeyNumpad8)
}

func (k keyboard) IsDownJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad2)
}

func (k keyboard) IsDownPressed() bool {
	return controller.IsRepeatingKeyboard(ebiten.KeyArrowDown, ebiten.KeyNumpad2)
}

func (k keyboard) IsLeftJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad4)
}

func (k keyboard) IsLeftPressed() bool {
	return controller.IsRepeatingKeyboard(ebiten.KeyArrowLeft, ebiten.KeyNumpad4)
}

func (k keyboard) IsRightJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad6)
}

func (k keyboard) IsRightPressed() bool {
	return controller.IsRepeatingKeyboard(ebiten.KeyArrowRight, ebiten.KeyNumpad6)
}

func (k keyboard) IsExitJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEscape)
}

func (k keyboard) IsExitPressed() bool {
	return controller.IsRepeatingKeyboard(ebiten.KeyEscape)
}

func (k keyboard) IsStartJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter)
}

func (k keyboard) IsStartPressed() bool {
	return controller.IsRepeatingKeyboard(ebiten.KeySpace, ebiten.KeyEnter, ebiten.KeyNumpadEnter)
}

func (k keyboard) Vibrate(_ time.Duration) {
	// nothing
}
