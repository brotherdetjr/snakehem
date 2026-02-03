package controller

import (
	"time"
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
