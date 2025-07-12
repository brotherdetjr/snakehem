package controller

import (
	"github.com/hajimehoshi/ebiten/v2"
	"snakehem/controller/gamepad"
	"snakehem/controller/keyboard"
)

func Controllers() []Controller {
	var result []Controller = nil
	result = append(result, keyboard.Instance)
	for _, g := range ebiten.AppendGamepadIDs(nil) {
		result = append(result, gamepad.NewGamepad(g))
	}
	return result
}
