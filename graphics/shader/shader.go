package shader

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
)

//go:embed crt_shader.kage
var shaderCode []byte

func NewShader() *ebiten.Shader {
	s, err := ebiten.NewShader(shaderCode)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	return s
}
