package game

import (
	"fmt"
	"math"
	"snakehem/graphics"
	"snakehem/graphics/pxterm24"
	"snakehem/graphics/shader"
	"snakehem/input/controller"
	"snakehem/model"
	"snakehem/model/perception"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pbnjay/pixfont"
	"github.com/rs/zerolog/log"
)

var scoreFmt = "%0" + fmt.Sprint(int(math.Log10(model.TargetScore))+1) + "d"
var pxterm16Height = pxterm24.Font.GetHeight()
var pxterm24Height = pxterm24.Font.GetHeight()

type Game struct {
	perception              perception.Perception
	controllers             []controller.Controller
	snakeControllers        []controller.Controller
	snakeHeadsRednessGrowth float32
	countdown               int
	elapsedFrames           uint64
	fadeCountdown           int
	applePresent            bool
	shader                  *ebiten.Shader
}

func Run() {
	pixfont.Spacing = 0
	// debug doesn't work well in fullscreen mode
	//ebiten.SetWindowSize(960, 960)
	ebiten.SetFullscreen(true)
	ebiten.SetTPS(model.Tps)
	ebiten.SetWindowTitle("snakehem")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	g := &Game{
		perception:              perception.NewPerception(),
		controllers:             nil,
		snakeControllers:        nil,
		snakeHeadsRednessGrowth: -1,
		countdown:               model.Tps * model.CountdownSeconds,
		elapsedFrames:           0,
		fadeCountdown:           0,
		applePresent:            false,
		shader:                  shader.NewShader(),
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return graphics.GridDimPx, graphics.GridDimPx
}
