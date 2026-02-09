package scoreboard

import (
	"fmt"
	"image/color"
	"snakehem/assets/pxterm16"
	"snakehem/assets/pxterm24"
	"snakehem/game/common"
	"snakehem/model"
	"snakehem/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

func (s *Scoreboard) Draw(screen *ebiten.Image) {
	vector.FillRect(
		screen,
		0,
		0,
		common.GridDimPx,
		common.GridDimPx,
		color.NRGBA{
			R: 85,
			G: 107,
			B: 47,
			A: 200,
		},
		false,
	)
	common.DrawTextCentered(
		screen,
		"GAME OVER",
		colornames.Yellow,
		float64(common.Pxterm24Height),
		pxterm24.Font,
	)
	common.DrawTextCentered(
		screen,
		"PRESS START BUTTON TO PLAY AGAIN",
		colornames.Yellow,
		float64(common.Pxterm24Height*2+common.Pxterm16Height),
		pxterm16.Font,
	)
	common.DrawTextCentered(
		screen,
		"      START                     ",
		color.White,
		float64(common.Pxterm24Height*2+common.Pxterm16Height),
		pxterm16.Font,
	)
	common.DrawTextCentered(
		screen,
		"OR SELECT BUTTON TO QUIT",
		colornames.Yellow,
		float64(common.Pxterm24Height*2+common.Pxterm16Height*2),
		pxterm16.Font,
	)
	common.DrawTextCentered(
		screen,
		"   SELECT               ",
		color.White,
		float64(common.Pxterm24Height*2+common.Pxterm16Height*2),
		pxterm16.Font,
	)
	for i, e := range s.entries {
		common.DrawTextCentered(
			screen,
			fmt.Sprintf("%s "+common.ScoreFmt, util.PadRight(e.Name, model.MaxNameLength), e.Score),
			e.ColourFunc(),
			float64(common.Pxterm24Height*2*(i+3)),
			pxterm24.Font,
		)
	}
}
