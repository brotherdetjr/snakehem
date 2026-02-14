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
	sharedContent     *shared.Content
	localContent      *local.Content
	unshadedContent   *unshaded.Content
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
		sharedContent:     shared.NewContent(),
		localContent:      local.NewContent(),
		unshadedContent:   unshaded.NewContent(),
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
