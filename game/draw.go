package game

import (
	"fmt"
	"image/color"
	"math"
	"slices"
	"snakehem/graphics"
	"snakehem/graphics/pxterm16"
	"snakehem/graphics/pxterm24"
	consts "snakehem/model"
	. "snakehem/model/apple"
	"snakehem/model/direction"
	. "snakehem/model/snake"
	. "snakehem/model/stage"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/pbnjay/pixfont"
	"golang.org/x/image/colornames"
)

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Darkolivegreen)
	g.drawItems(screen)
	switch g.stage {
	case Lobby:
		g.drawScores(screen)
		snakeCount := len(g.snakes)
		if snakeCount < 2 {
			drawTextCentered(
				screen,
				"PLAYERS PRESS ANY BUTTON TO JOIN",
				colornames.Yellow,
				graphics.GridDimPx/2.5,
				pxterm16.Font,
			)
		} else {
			drawTextCentered(
				screen,
				"PLAYERS PRESS START BUTTON TO GO",
				colornames.Yellow,
				graphics.GridDimPx/2.5,
				pxterm16.Font,
			)
			drawTextCentered(
				screen,
				"              START             ",
				color.White,
				graphics.GridDimPx/2.5,
				pxterm16.Font,
			)
			if snakeCount < consts.MaxSnakes {
				drawTextCentered(
					screen,
					"OR ANY OTHER BUTTON TO JOIN",
					colornames.Yellow,
					graphics.GridDimPx/2.5+float64(pxterm16Height)*1.5,
					pxterm16.Font,
				)
			}
		}
	case Action:
		if g.fadeCountdown > 0 {
			vector.DrawFilledRect(
				screen,
				0,
				0,
				graphics.GridDimPx,
				graphics.GridDimPx,
				color.NRGBA{
					R: 85,
					G: 107,
					B: 47,
					A: uint8((consts.GridFadeCountdown - g.fadeCountdown) * 200 / consts.GridFadeCountdown),
				},
				false,
			)
		}
		g.drawScores(screen)
		g.drawCountdown(screen)
		g.drawTimeElapsed(screen)
	case Scoreboard:
		g.drawScoreboard(screen)
		g.drawTimeElapsed(screen)
	}
	g.applyShader(screen)
}

func (g *Game) drawScoreboard(screen *ebiten.Image) {
	vector.DrawFilledRect(
		screen,
		0,
		0,
		graphics.GridDimPx,
		graphics.GridDimPx,
		color.NRGBA{
			R: 85,
			G: 107,
			B: 47,
			A: 200,
		},
		false,
	)
	drawTextCentered(
		screen,
		"GAME OVER",
		colornames.Yellow,
		float64(pxterm24Height),
		pxterm24.Font,
	)
	drawTextCentered(
		screen,
		"PRESS START BUTTON TO PLAY AGAIN",
		colornames.Yellow,
		float64(pxterm24Height*2+pxterm16Height),
		pxterm16.Font,
	)
	drawTextCentered(
		screen,
		"      START                     ",
		color.White,
		float64(pxterm24Height*2+pxterm16Height),
		pxterm16.Font,
	)
	drawTextCentered(
		screen,
		"OR SELECT BUTTON TO QUIT",
		colornames.Yellow,
		float64(pxterm24Height*2+pxterm16Height*2),
		pxterm16.Font,
	)
	drawTextCentered(
		screen,
		"   SELECT               ",
		color.White,
		float64(pxterm24Height*2+pxterm16Height*2),
		pxterm16.Font,
	)
	snakes := make([]*Snake, len(g.snakes))
	copy(snakes, g.snakes)
	slices.SortFunc(snakes, func(a, b *Snake) int {
		return b.Score - a.Score
	})
	for i, snake := range snakes {
		top := pxterm24Height * 2 * (i + 3)
		score := snake.Score
		if score > consts.TargetScore {
			score = consts.TargetScore
		}
		drawTextCentered(
			screen,
			fmt.Sprintf("PLAYER %d "+scoreFmt, i+1, score),
			withRedness(snake.Colour, snake.Links[0].Redness),
			float64(top),
			pxterm24.Font,
		)
	}
}

func (g *Game) drawTimeElapsed(screen *ebiten.Image) {
	t := time.UnixMilli(int64(float32(g.elapsedFrames) / consts.Tps * 1000))
	drawTextCentered(
		screen,
		t.Format("04:05.0"),
		colornames.White,
		graphics.GridDimPx-float64(pxterm16Height)*1.5,
		pxterm16.Font,
	)
}

func (g *Game) drawCountdown(screen *ebiten.Image) {
	if g.countdown <= 0 {
		return
	}
	count := int(math.Ceil(float64(g.countdown)/consts.Tps)) - 1
	var txt string
	switch count {
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
	drawTextCentered(screen, txt, color.White, graphics.GridDimPx/2.5, pxterm24.Font)
	if count > 0 {
		drawTextCentered(
			screen,
			fmt.Sprintf("TARGET SCORE: %d", consts.TargetScore),
			colornames.Yellow,
			graphics.GridDimPx/2.5+float64(pxterm24Height*2),
			pxterm24.Font,
		)
	}
}

func drawTextCentered(screen *ebiten.Image, txt string, colour color.Color, top float64, font *pixfont.PixFont) {
	txtWidth := font.MeasureString(txt)
	font.DrawString(screen, (graphics.GridDimPx-txtWidth)/2, int(top), txt, colour)
}

func (g *Game) drawItems(screen *ebiten.Image) {
	for i := 0; i < consts.GridSize; i++ {
		for j := 0; j < consts.GridSize; j++ {
			if val := g.grid[i][j]; val != nil {
				switch item := val.(type) {
				case *Link:
					shrink := (1 - float32(item.HealthPercent)/100) * graphics.CellDimPx * 0.5
					if item != item.Snake.Links[0] {
						vector.DrawFilledRect(
							screen,
							float32(item.X*graphics.CellDimPx)+shrink,
							float32(item.Y*graphics.CellDimPx)+shrink,
							graphics.CellDimPx-shrink*2,
							graphics.CellDimPx-shrink*2,
							withRedness(item.Snake.Colour, item.Redness),
							false,
						)
					} else if g.countdown > consts.Tps {
						vector.DrawFilledRect(
							screen,
							float32(item.X*graphics.CellDimPx)+shrink,
							float32(item.Y*graphics.CellDimPx)+shrink,
							graphics.CellDimPx-shrink*2,
							graphics.CellDimPx-shrink*2,
							withRedness(item.Snake.Colour, item.Redness),
							false,
						)
					} else {
						var x1, y1, x2, y2 float32
						switch item.Snake.Direction {
						case direction.Up:
							x1 = float32(item.X*graphics.CellDimPx) + graphics.EyeGapPx
							y1 = float32(item.Y*graphics.CellDimPx) + graphics.EyeGapPx
							x2 = float32((item.X+1)*graphics.CellDimPx) - graphics.EyeGapPx
							y2 = float32(item.Y*graphics.CellDimPx) + graphics.EyeGapPx
						case direction.Down:
							x1 = float32(item.X*graphics.CellDimPx) + graphics.EyeGapPx
							y1 = float32((item.Y+1)*graphics.CellDimPx) - graphics.EyeGapPx
							x2 = float32((item.X+1)*graphics.CellDimPx) - graphics.EyeGapPx
							y2 = float32((item.Y+1)*graphics.CellDimPx) - graphics.EyeGapPx
						case direction.Left:
							x1 = float32(item.X*graphics.CellDimPx) + graphics.EyeGapPx
							y1 = float32((item.Y+1)*graphics.CellDimPx) - graphics.EyeGapPx
							x2 = float32(item.X*graphics.CellDimPx) + graphics.EyeGapPx
							y2 = float32(item.Y*graphics.CellDimPx) + graphics.EyeGapPx
						case direction.Right:
							x1 = float32((item.X+1)*graphics.CellDimPx) - graphics.EyeGapPx
							y1 = float32((item.Y+1)*graphics.CellDimPx) - graphics.EyeGapPx
							x2 = float32((item.X+1)*graphics.CellDimPx) - graphics.EyeGapPx
							y2 = float32(item.Y*graphics.CellDimPx) + graphics.EyeGapPx
						case direction.None:
						}
						if x1 != 0 || y1 != 0 || x2 != 0 || y2 != 0 {
							vector.DrawFilledCircle(
								screen,
								x1,
								y1,
								graphics.EyeRadiusPx,
								withRedness(item.Snake.Colour, item.Redness),
								false,
							)
							vector.DrawFilledCircle(
								screen,
								x2,
								y2,
								graphics.EyeRadiusPx,
								withRedness(item.Snake.Colour, item.Redness),
								false,
							)
						} else {
							vector.DrawFilledRect(
								screen,
								float32(item.X*graphics.CellDimPx),
								float32(item.Y*graphics.CellDimPx),
								graphics.CellDimPx,
								graphics.CellDimPx,
								withRedness(item.Snake.Colour, item.Redness),
								false,
							)
						}
					}
				case *Apple:
					vector.DrawFilledRect(
						screen,
						float32(item.X*graphics.CellDimPx),
						float32(item.Y*graphics.CellDimPx),
						graphics.CellDimPx,
						graphics.CellDimPx,
						colornames.Red,
						false,
					)
				}
			}
		}
	}
}

func (g *Game) drawScores(screen *ebiten.Image) {
	scoresAtTop := len(g.snakes)
	if scoresAtTop > graphics.MaxScoresAtTop {
		scoresAtTop = graphics.MaxScoresAtTop
	}
	g.drawScoreRow(screen, g.snakes[:scoresAtTop], pxterm24Height/2)
	// when there are many players, not all scores can be fit in one line
	g.drawScoreRow(screen, g.snakes[scoresAtTop:], graphics.GridDimPx-pxterm24Height-pxterm16Height*2)
}

func (g *Game) drawScoreRow(screen *ebiten.Image, snakes []*Snake, rowTopPos int) {
	span := float64(screen.Bounds().Dx()) / float64(len(snakes))
	for i, snake := range snakes {
		if g.stage != Action || snake.Score+consts.ApproachingTargetScoreGap < consts.TargetScore || (g.elapsedFrames/(consts.Tps/4))%2 > 0 {
			txt, colour := g.scoreStrAndColourForIthSnake(snake)
			x := int(span*float64(i) + span/2 - float64(pxterm24.Font.MeasureString(txt))/2 + 2)
			pxterm24.Font.DrawString(screen, x, rowTopPos, txt, colour)
		}
	}
}

func (g *Game) applyShader(screen *ebiten.Image) {
	w := screen.Bounds().Dx()
	h := screen.Bounds().Dy()
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = screen
	opts.Uniforms = map[string]interface{}{
		// Kage uniforms here
	}
	img := ebiten.NewImage(w, h)
	img.DrawRectShader(w, h, g.shader, opts)
	screen.DrawImage(img, nil)
}

func (g *Game) scoreStrAndColourForIthSnake(snake *Snake) (string, color.Color) {
	score := snake.Score
	if score > consts.TargetScore {
		score = consts.TargetScore
	}
	txt := fmt.Sprintf(scoreFmt, score)
	var colour color.Color
	if g.stage == Action && g.countdown <= consts.Tps {
		colour = snake.Colour
	} else {
		colour = withRedness(snake.Colour, snake.Links[0].Redness)
	}
	return txt, colour
}

func withRedness(colour color.Color, redness float32) color.Color {
	red, green, blue, _ := colour.RGBA()
	r := float32(red >> 8)
	g := float32(green >> 8)
	b := float32(blue >> 8)
	return color.NRGBA{
		R: uint8(r + (255-r)*redness),
		G: uint8(g * (1 - redness)),
		B: uint8(b * (1 - redness)),
		A: 255,
	}
}
