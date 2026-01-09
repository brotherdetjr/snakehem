package controllers

import (
	"github.com/hajimehoshi/ebiten/v2"
	"snakehem/controllers/controller"
	"snakehem/controllers/gamepad"
	"snakehem/controllers/keyboard"
	"snakehem/controllers/keyboardwasd"
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
