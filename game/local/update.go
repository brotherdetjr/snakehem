package local

import (
	"snakehem/game/common"
)

func (c *Content) Update(ctx *common.Context) {
	if c.textInput != nil {
		c.textInput.Update(ctx)
	}
}
