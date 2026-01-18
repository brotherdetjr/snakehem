package controller

import (
	"time"
)

type Controller interface {
	Equals(controller Controller) bool
	IsAnyJustPressed() bool
	IsUpJustPressed() bool
	IsDownJustPressed() bool
	IsLeftJustPressed() bool
	IsRightJustPressed() bool
	IsExitJustPressed() bool
	IsStartJustPressed() bool
	Vibrate(duration time.Duration)
}
