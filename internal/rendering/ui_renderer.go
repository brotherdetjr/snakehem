package rendering

import (
	"fmt"
	"image/color"
	"math"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/pbnjay/pixfont"
	"golang.org/x/image/colornames"

	"snakehem/internal/config"
	"snakehem/internal/entities"
)

// UIRenderer handles all UI text and overlay drawing
type UIRenderer struct {
	config      *config.GameConfig
	fontManager *FontManager
}

// NewUIRenderer creates a new UI renderer
func NewUIRenderer(config *config.GameConfig, fontManager *FontManager) *UIRenderer {
	return &UIRenderer{
		config:      config,
		fontManager: fontManager,
	}
}

// DrawLobbyUI draws the lobby screen instructions
func (ur *UIRenderer) DrawLobbyUI(screen *ebiten.Image, snakeCount int) {
	font16 := ur.fontManager.GetFont16()
	font16Height := font16.GetHeight()
	gridDimPx := float64(ur.config.GridDimPx())

	if snakeCount < 2 {
		ur.drawTextCentered(
			screen,
			"PLAYERS PRESS ANY BUTTON TO JOIN",
			colornames.Yellow,
			gridDimPx/2.5,
			font16,
		)
	} else {
		ur.drawTextCentered(
			screen,
			"PLAYERS PRESS START BUTTON TO GO",
			colornames.Yellow,
			gridDimPx/2.5,
			font16,
		)
		ur.drawTextCentered(
			screen,
			"              START             ",
			color.White,
			gridDimPx/2.5,
			font16,
		)
		if snakeCount < ur.config.MaxSnakes {
			ur.drawTextCentered(
				screen,
				"OR ANY OTHER BUTTON TO JOIN",
				colornames.Yellow,
				gridDimPx/2.5+float64(font16Height)*1.5,
				font16,
			)
		}
	}
}

// DrawCountdown draws the countdown timer and target score
func (ur *UIRenderer) DrawCountdown(screen *ebiten.Image, countdown int, elapsedFrames uint64) {
	if countdown <= 0 {
		return
	}

	font24 := ur.fontManager.GetFont24()
	font24Height := font24.GetHeight()
	gridDimPx := float64(ur.config.GridDimPx())

	count := int(math.Ceil(float64(countdown)/float64(ur.config.Tps()))) - 1
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

	ur.drawTextCentered(screen, txt, color.White, gridDimPx/2.5, font24)

	if count > 0 {
		ur.drawTextCentered(
			screen,
			fmt.Sprintf("TARGET SCORE: %d", ur.config.TargetScore),
			colornames.Yellow,
			gridDimPx/2.5+float64(font24Height*2),
			font24,
		)
	}
}

// DrawScoreboard draws the game over screen with final scores
func (ur *UIRenderer) DrawScoreboard(screen *ebiten.Image, snakes []*entities.Snake) {
	font16 := ur.fontManager.GetFont16()
	font24 := ur.fontManager.GetFont24()
	font16Height := font16.GetHeight()
	font24Height := font24.GetHeight()
	gridDimPx := ur.config.GridDimPx()

	// Draw semi-transparent overlay
	vector.DrawFilledRect(
		screen,
		0,
		0,
		float32(gridDimPx),
		float32(gridDimPx),
		color.NRGBA{
			R: 85,
			G: 107,
			B: 47,
			A: 200,
		},
		false,
	)

	// Draw title
	ur.drawTextCentered(
		screen,
		"GAME OVER",
		colornames.Yellow,
		float64(font24Height),
		font24,
	)

	// Draw restart instructions
	ur.drawTextCentered(
		screen,
		"PRESS START BUTTON TO PLAY AGAIN",
		colornames.Yellow,
		float64(font24Height*2+font16Height),
		font16,
	)
	ur.drawTextCentered(
		screen,
		"      START                     ",
		color.White,
		float64(font24Height*2+font16Height),
		font16,
	)

	// Draw quit instructions
	ur.drawTextCentered(
		screen,
		"OR SELECT BUTTON TO QUIT",
		colornames.Yellow,
		float64(font24Height*2+font16Height*2),
		font16,
	)
	ur.drawTextCentered(
		screen,
		"   SELECT               ",
		color.White,
		float64(font24Height*2+font16Height*2),
		font16,
	)

	// Draw sorted scores
	sortedSnakes := make([]*entities.Snake, len(snakes))
	copy(sortedSnakes, snakes)
	slices.SortFunc(sortedSnakes, func(a, b *entities.Snake) int {
		return b.Score - a.Score
	})

	scoreFmt := ur.getScoreFormat()
	for i, snake := range sortedSnakes {
		top := font24Height * 2 * (i + 3)
		score := snake.Score
		if score > ur.config.TargetScore {
			score = ur.config.TargetScore
		}
		ur.drawTextCentered(
			screen,
			fmt.Sprintf("PLAYER %d "+scoreFmt, i+1, score),
			withRedness(snake.Colour, snake.Links[0].Redness),
			float64(top),
			font24,
		)
	}
}

// DrawScores draws scores at top/bottom of screen during gameplay
func (ur *UIRenderer) DrawScores(screen *ebiten.Image, snakes []*entities.Snake, elapsedFrames uint64, isActionState bool, countdown int) {
	font24 := ur.fontManager.GetFont24()
	font16 := ur.fontManager.GetFont16()
	font24Height := font24.GetHeight()
	font16Height := font16.GetHeight()
	gridDimPx := ur.config.GridDimPx()

	scoresAtTop := len(snakes)
	if scoresAtTop > ur.config.MaxScoresAtTop {
		scoresAtTop = ur.config.MaxScoresAtTop
	}

	ur.drawScoreRow(screen, snakes[:scoresAtTop], font24Height/2, elapsedFrames, isActionState, countdown)
	ur.drawScoreRow(screen, snakes[scoresAtTop:], gridDimPx-font24Height-font16Height*2, elapsedFrames, isActionState, countdown)
}

// drawScoreRow draws a row of scores
func (ur *UIRenderer) drawScoreRow(screen *ebiten.Image, snakes []*entities.Snake, rowTopPos int, elapsedFrames uint64, isActionState bool, countdown int) {
	if len(snakes) == 0 {
		return
	}

	font24 := ur.fontManager.GetFont24()
	span := float64(screen.Bounds().Dx()) / float64(len(snakes))
	scoreFmt := ur.getScoreFormat()

	for i, snake := range snakes {
		// Blinking effect when approaching target score
		shouldBlink := isActionState && snake.Score+ur.config.ApproachingTargetScoreGap >= ur.config.TargetScore
		isBlinkVisible := (elapsedFrames/(uint64(ur.config.Tps())/4))%2 > 0

		if !shouldBlink || isBlinkVisible {
			score := snake.Score
			if score > ur.config.TargetScore {
				score = ur.config.TargetScore
			}
			txt := fmt.Sprintf(scoreFmt, score)

			// Determine color based on state
			var colour color.Color
			if isActionState && countdown <= ur.config.Tps() {
				colour = snake.Colour
			} else {
				colour = withRedness(snake.Colour, snake.Links[0].Redness)
			}

			x := int(span*float64(i) + span/2 - float64(font24.MeasureString(txt))/2 + 2)
			font24.DrawString(screen, x, rowTopPos, txt, colour)
		}
	}
}

// DrawTimeElapsed draws the elapsed time at the bottom of the screen
func (ur *UIRenderer) DrawTimeElapsed(screen *ebiten.Image, elapsedFrames uint64) {
	font16 := ur.fontManager.GetFont16()
	font16Height := font16.GetHeight()
	gridDimPx := float64(ur.config.GridDimPx())

	t := time.UnixMilli(int64(float32(elapsedFrames) / float32(ur.config.Tps()) * 1000))
	ur.drawTextCentered(
		screen,
		t.Format("04:05.0"),
		colornames.White,
		gridDimPx-float64(font16Height)*1.5,
		font16,
	)
}

// DrawFadeOverlay draws the fade overlay when someone reaches the target score
func (ur *UIRenderer) DrawFadeOverlay(screen *ebiten.Image, fadeCountdown int) {
	if fadeCountdown <= 0 {
		return
	}

	gridDimPx := float32(ur.config.GridDimPx())
	fadeProgress := ur.config.GridFadeCountdown - fadeCountdown
	alpha := uint8(fadeProgress * 200 / ur.config.GridFadeCountdown)

	vector.DrawFilledRect(
		screen,
		0,
		0,
		gridDimPx,
		gridDimPx,
		color.NRGBA{
			R: 85,
			G: 107,
			B: 47,
			A: alpha,
		},
		false,
	)
}

// drawTextCentered draws text centered horizontally at the given vertical position
func (ur *UIRenderer) drawTextCentered(screen *ebiten.Image, txt string, colour color.Color, top float64, font *pixfont.PixFont) {
	txtWidth := font.MeasureString(txt)
	x := (ur.config.GridDimPx() - txtWidth) / 2
	font.DrawString(screen, x, int(top), txt, colour)
}

// getScoreFormat returns the format string for scores
func (ur *UIRenderer) getScoreFormat() string {
	return fmt.Sprintf("%%0%dd", int(math.Log10(float64(ur.config.TargetScore)))+1)
}
