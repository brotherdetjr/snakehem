package input

import (
	"snakehem/input/controller"
	"snakehem/input/gamepad"
	"snakehem/input/keyboard"
	"snakehem/input/keyboardwasd"

	"github.com/hajimehoshi/ebiten/v2"
)

func Controllers() []controller.Controller {
	var result []controller.Controller = nil
	result = append(result, keyboard.Instance)
	result = append(result, keyboardwasd.Instance)
	for _, g := range ebiten.AppendGamepadIDs(nil) {
		result = append(result, gamepad.NewGamepad(g))
	}
	return result
}
