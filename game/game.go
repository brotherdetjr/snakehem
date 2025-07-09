package game

import (
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pbnjay/pixfont"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	. "snakehem/controller"
	"snakehem/pxterm24"
	. "snakehem/snake"
	. "snakehem/state"
)

const (
	tps                           = 10
	gridSize                      = 63
	cellDimPx                     = 11
	gridDimPx                     = gridSize * cellDimPx
	maxScoresAtTop                = 5
	countdownSeconds              = 4
	snakeTargetLength             = 50
	healthReductionPerBite        = 10
	nippedTailLinkBonusMultiplier = 2
	bitLinkScore                  = 1
	appleScore                    = 50
	targetScore                   = 999
	maxSnakes                     = len(snakeColours)
	approachingTargetScoreGap     = snakeTargetLength * nippedTailLinkBonusMultiplier * 1.2
	gridFadeCountdown             = 15
	newAppleProbabilityParam      = tps * 3
)

//go:embed crt_shader.kage
var shaderCode []byte
var shader = newShader()
var scoreFmt = "%0" + fmt.Sprint(int(math.Log10(targetScore))+1) + "d"
var pxterm16Height = pxterm24.Font.GetHeight()
var pxterm24Height = pxterm24.Font.GetHeight()
var snakeColours = [...]color.Color{
	colornames.Lightgrey,
	color.NRGBA{
		R: 255,
		G: 128,
		B: 10,
		A: 255,
	},
	colornames.Yellow,
	color.NRGBA{
		R: 100,
		G: 170,
		B: 0,
		A: 255,
	},
	colornames.Cyan,
	colornames.Blue,
	color.NRGBA{
		R: 0,
		G: 0,
		B: 100,
		A: 255,
	},
	color.NRGBA{
		R: 100,
		G: 0,
		B: 84,
		A: 255,
	},
	colornames.Magenta,
}

type Game struct {
	grid          [gridSize][gridSize]any
	snakes        []*Snake
	controllers   []Controller
	state         State
	countdown     int
	elapsedFrames uint64
	fadeCountdown int
	applePresent  bool
}

func Run() {
	pixfont.Spacing = 0
	// debug doesn't work well in fullscreen mode
	ebiten.SetWindowSize(960, 960)
	//ebiten.SetFullscreen(true)
	ebiten.SetTPS(tps)
	ebiten.SetWindowTitle("snakehem")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	g := &Game{
		grid:          [gridSize][gridSize]any{},
		snakes:        nil,
		controllers:   nil,
		state:         Lobby,
		countdown:     tps * countdownSeconds,
		elapsedFrames: 0,
		fadeCountdown: 0,
		applePresent:  false,
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return gridDimPx, gridDimPx
}

func newShader() *ebiten.Shader {
	s, err := ebiten.NewShader(shaderCode)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	return s
}
