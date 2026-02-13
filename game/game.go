package game

import (
	"snakehem/assets/shader"
	"snakehem/game/common"
	"snakehem/game/local"
	"snakehem/game/shared"
	"snakehem/game/unshaded"
	"snakehem/input/controller"
	"snakehem/model"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pbnjay/pixfont"
	"github.com/rs/zerolog/log"
)

type Game struct {
	sharedState       *shared.State
	localState        *local.State
	unshadedState     *unshaded.State
	controllers       []controller.Controller
	activeControllers []controller.Controller
	shader            *ebiten.Shader
}

func Run() {
	pixfont.Spacing = 0
	// debug doesn't work well in fullscreen mode
	//ebiten.SetWindowSize(960, 960)
	ebiten.SetFullscreen(true)
	ebiten.SetTPS(model.Tps)
	ebiten.SetWindowTitle("snakehem")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetScreenClearedEveryFrame(false)
	g := &Game{
		sharedState:       shared.NewSharedState(),
		localState:        local.NewLocalState(),
		unshadedState:     unshaded.NewUnshadedState(),
		controllers:       nil,
		activeControllers: nil,
		shader:            shader.NewShader(),
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return common.GridDimPx, common.GridDimPx
}
