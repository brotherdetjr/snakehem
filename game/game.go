package game

import (
	"snakehem/assets/shader"
	"snakehem/game/common"
	"snakehem/game/local"
	"snakehem/game/shared"
	"snakehem/input/controller"
	"snakehem/model"
	"snakehem/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pbnjay/pixfont"
	"github.com/rs/zerolog/log"
)

type Game struct {
	sharedState        *shared.State
	localState         *local.State
	controllers        []controller.Controller
	activeControllers  []controller.Controller
	shader             *ebiten.Shader
	lastFrame          *ebiten.Image
	perfTracker        *util.PerfTracker
	perfTrackerVisible bool
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
		sharedState:       shared.NewSharedState(),
		localState:        local.NewLocalState(),
		controllers:       nil,
		activeControllers: nil,
		shader:            shader.NewShader(),
		lastFrame:         ebiten.NewImage(common.GridDimPx, common.GridDimPx),
		perfTracker:       nil,
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return common.GridDimPx, common.GridDimPx
}
