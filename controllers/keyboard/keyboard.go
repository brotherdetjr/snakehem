package keyboard

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"snakehem/controllers/controller"
	"time"
)

type keyboard struct {
}

var Instance controller.Controller = keyboard{}

func (k keyboard) Equals(controller controller.Controller) bool {
	_, ok := controller.(keyboard)
	return ok
}

func (k keyboard) IsAnyJustPressed() bool {
	return len(inpututil.AppendJustPressedKeys(nil)) > 0
}

func (k keyboard) IsUpJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad8)
}

func (k keyboard) IsDownJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad2)
}

func (k keyboard) IsLeftJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad4)
}

func (k keyboard) IsRightJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad6)
}

func (k keyboard) IsExitJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEscape)
}

func (k keyboard) IsStartJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter)
}

func (k keyboard) Vibrate(_ time.Duration) {
	// nothing
}
