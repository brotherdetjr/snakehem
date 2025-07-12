package game

import (
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pbnjay/pixfont"
	"github.com/rs/zerolog/log"
	"math"
	"snakehem/consts"
	"snakehem/controllers/controller"
	"snakehem/pxterm24"
	. "snakehem/snake"
	. "snakehem/state"
)

//go:embed crt_shader.kage
var shaderCode []byte
var shader = newShader()
var scoreFmt = "%0" + fmt.Sprint(int(math.Log10(consts.TargetScore))+1) + "d"
var pxterm16Height = pxterm24.Font.GetHeight()
var pxterm24Height = pxterm24.Font.GetHeight()

type Game struct {
	grid          [consts.GridSize][consts.GridSize]any
	snakes        []*Snake
	controllers   []controller.Controller
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
	ebiten.SetTPS(consts.Tps)
	ebiten.SetWindowTitle("snakehem")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	g := &Game{
		grid:          [consts.GridSize][consts.GridSize]any{},
		snakes:        nil,
		controllers:   nil,
		state:         Lobby,
		countdown:     consts.Tps * consts.CountdownSeconds,
		elapsedFrames: 0,
		fadeCountdown: 0,
		applePresent:  false,
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return consts.GridDimPx, consts.GridDimPx
}

func newShader() *ebiten.Shader {
	s, err := ebiten.NewShader(shaderCode)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	return s
}
