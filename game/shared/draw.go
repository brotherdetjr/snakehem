package shared

import (
	"fmt"
	"image/color"
	"slices"
	"snakehem/assets/pxterm16"
	"snakehem/assets/pxterm24"
	"snakehem/game/common"
	"snakehem/model"
	"snakehem/model/apple"
	"snakehem/model/direction"
	"snakehem/model/snake"
	"snakehem/util"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

const (
	MaxScoresAtTop = 5
	EyeRadiusPx    = 2
	EyeGapPx       = 3
)

func (p *State) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Darkolivegreen)
	drawItems(p, screen)
	switch p.Stage {
	case Lobby:
		drawScores(p, screen)
		snakeCount := len(p.Snakes)
		if snakeCount < 2 {
			common.DrawTextCentered(
				screen,
				"PLAYERS PRESS ANY BUTTON TO JOIN",
				colornames.Yellow,
				common.GridDimPx/2.5,
				pxterm16.Font,
			)
		} else {
			common.DrawTextCentered(
				screen,
				"PLAYERS PRESS START BUTTON TO GO",
				colornames.Yellow,
				common.GridDimPx/2.5,
				pxterm16.Font,
			)
			common.DrawTextCentered(
				screen,
				"              START             ",
				color.White,
				common.GridDimPx/2.5,
				pxterm16.Font,
			)
			if snakeCount < model.MaxSnakes {
				common.DrawTextCentered(
					screen,
					"OR ANY OTHER BUTTON TO JOIN",
					colornames.Yellow,
					common.GridDimPx/2.5+float64(common.Pxterm16Height)*1.5,
					pxterm16.Font,
				)
			}
		}
	case Action:
		if p.FadeCountdown > 0 {
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
					A: uint8((model.GridFadeCountdown - p.FadeCountdown) * 200 / model.GridFadeCountdown),
				},
				false,
			)
		}
		drawScores(p, screen)
		drawCountdown(p, screen)
		drawTimeElapsed(p, screen)
	case Scoreboard:
		drawScoreboard(p, screen)
		drawTimeElapsed(p, screen)
	}
}

func drawItems(p *State, screen *ebiten.Image) {
	for i := 0; i < model.GridSize; i++ {
		for j := 0; j < model.GridSize; j++ {
			if val := p.Grid[i][j]; val != nil {
				switch item := val.(type) {
				case *snake.Link:
					s := p.Snakes[item.SnakeId]
					shrink := (1 - float32(item.HealthPercent)/100) * common.CellDimPx * 0.5
					if item != s.Links[0] || p.Countdown > 0 {
						vector.FillRect(
							screen,
							float32(item.X*common.CellDimPx)+shrink,
							float32(item.Y*common.CellDimPx)+shrink,
							common.CellDimPx-shrink*2,
							common.CellDimPx-shrink*2,
							common.WithRedness(s.Colour, item.Redness),
							false,
						)
					} else {
						var x1, y1, x2, y2 float32
						switch s.Direction {
						case direction.Up:
							x1 = float32(item.X*common.CellDimPx) + EyeGapPx
							y1 = float32(item.Y*common.CellDimPx) + EyeGapPx
							x2 = float32((item.X+1)*common.CellDimPx) - EyeGapPx
							y2 = float32(item.Y*common.CellDimPx) + EyeGapPx
						case direction.Down:
							x1 = float32(item.X*common.CellDimPx) + EyeGapPx
							y1 = float32((item.Y+1)*common.CellDimPx) - EyeGapPx
							x2 = float32((item.X+1)*common.CellDimPx) - EyeGapPx
							y2 = float32((item.Y+1)*common.CellDimPx) - EyeGapPx
						case direction.Left:
							x1 = float32(item.X*common.CellDimPx) + EyeGapPx
							y1 = float32((item.Y+1)*common.CellDimPx) - EyeGapPx
							x2 = float32(item.X*common.CellDimPx) + EyeGapPx
							y2 = float32(item.Y*common.CellDimPx) + EyeGapPx
						case direction.Right:
							x1 = float32((item.X+1)*common.CellDimPx) - EyeGapPx
							y1 = float32((item.Y+1)*common.CellDimPx) - EyeGapPx
							x2 = float32((item.X+1)*common.CellDimPx) - EyeGapPx
							y2 = float32(item.Y*common.CellDimPx) + EyeGapPx
						case direction.None:
						}
						if x1 != 0 || y1 != 0 || x2 != 0 || y2 != 0 {
							vector.FillCircle(
								screen,
								x1,
								y1,
								EyeRadiusPx,
								common.WithRedness(s.Colour, item.Redness),
								false,
							)
							vector.FillCircle(
								screen,
								x2,
								y2,
								EyeRadiusPx,
								common.WithRedness(s.Colour, item.Redness),
								false,
							)
						} else {
							vector.FillRect(
								screen,
								float32(item.X*common.CellDimPx),
								float32(item.Y*common.CellDimPx),
								common.CellDimPx,
								common.CellDimPx,
								common.WithRedness(s.Colour, item.Redness),
								false,
							)
						}
					}
				case *apple.Apple:
					vector.FillRect(
						screen,
						float32(item.X*common.CellDimPx),
						float32(item.Y*common.CellDimPx),
						common.CellDimPx,
						common.CellDimPx,
						colornames.Red,
						false,
					)
				}
			}
		}
	}
}

func drawScores(p *State, screen *ebiten.Image) {
	snakes := p.Snakes
	scoresAtTop := len(snakes)
	if scoresAtTop > MaxScoresAtTop {
		scoresAtTop = MaxScoresAtTop
	}
	drawScoreRow(p, screen, snakes[:scoresAtTop], common.Pxterm24Height/2)
	// when there are many players, not all scores can be fit in one line
	drawScoreRow(p, screen, snakes[scoresAtTop:], common.GridDimPx-common.Pxterm24Height-common.Pxterm16Height*2)
}

func drawScoreRow(p *State, screen *ebiten.Image, snakes []*snake.Snake, rowTopPos int) {
	span := float64(screen.Bounds().Dx()) / float64(len(snakes))
	for i, s := range snakes {
		if p.Stage != Action || s.Score+model.ApproachingTargetScoreGap < model.TargetScore || (p.ActionFrameCount/(model.Tps/4))%2 > 0 {
			txt, colour := scoreStrAndColourForIthSnake(p, s)
			x := int(span*float64(i) + span/2 - float64(pxterm24.Font.MeasureString(txt))/2 + 2)
			pxterm24.Font.DrawString(screen, x, rowTopPos, txt, colour)
		}
	}
}

func scoreStrAndColourForIthSnake(p *State, snake *snake.Snake) (string, color.Color) {
	score := snake.Score
	if score > model.TargetScore {
		score = model.TargetScore
	}
	txt := fmt.Sprintf(common.ScoreFmt, score)
	var colour color.Color
	if p.Stage == Action && p.Countdown < 1 {
		colour = snake.Colour
	} else {
		colour = common.WithRedness(snake.Colour, snake.Links[0].Redness)
	}
	return txt, colour
}

func drawScoreboard(p *State, screen *ebiten.Image) {
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
	snakes := make([]*snake.Snake, len(p.Snakes))
	copy(snakes, p.Snakes)
	slices.SortFunc(snakes, func(a, b *snake.Snake) int {
		return b.Score - a.Score
	})
	for i, s := range snakes {
		top := common.Pxterm24Height * 2 * (i + 3)
		score := s.Score
		if score > model.TargetScore {
			score = model.TargetScore
		}
		common.DrawTextCentered(
			screen,
			fmt.Sprintf("%s "+common.ScoreFmt, util.PadRight(s.Name, model.MaxNameLength), score),
			common.WithRedness(s.Colour, s.Links[0].Redness),
			float64(top),
			pxterm24.Font,
		)
	}
}

func drawTimeElapsed(p *State, screen *ebiten.Image) {
	t := time.UnixMilli(int64(float32(p.ActionFrameCount) / model.Tps * 1000))
	common.DrawTextCentered(
		screen,
		t.Format("04:05.0"),
		colornames.White,
		common.GridDimPx-float64(common.Pxterm16Height)*1.5,
		pxterm16.Font,
	)
}

func drawCountdown(p *State, screen *ebiten.Image) {
	if p.Countdown <= 0 {
		return
	}
	var txt string
	switch p.Countdown {
	case 3:
		txt = "THREE"
	case 2:
		txt = "TWO"
	case 1:
		txt = "ONE"
	case 0:
		txt = "GO!"
	default:
		txt = "WAIT..."
	}
	common.DrawTextCentered(screen, txt, color.White, common.GridDimPx/2.5, pxterm24.Font)
	if p.Countdown > 0 {
		common.DrawTextCentered(
			screen,
			fmt.Sprintf("TARGET SCORE: %d", model.TargetScore),
			colornames.Yellow,
			common.GridDimPx/2.5+float64(common.Pxterm24Height*2),
			pxterm24.Font,
		)
	}
}
