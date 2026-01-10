package rendering

import (
	_ "embed"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed crt_shader.kage
var shaderCode []byte

// PostProcessor handles CRT shader post-processing effects
type PostProcessor struct {
	shader *ebiten.Shader
}

// NewPostProcessor creates a post-processor with the embedded CRT shader
// Returns error instead of using log.Fatal
func NewPostProcessor() (*PostProcessor, error) {
	return NewPostProcessorWithShader(shaderCode)
}

// NewPostProcessorWithShader creates a post-processor with custom shader code
// Returns error instead of using log.Fatal
func NewPostProcessorWithShader(code []byte) (*PostProcessor, error) {
	shader, err := ebiten.NewShader(code)
	if err != nil {
		return nil, fmt.Errorf("failed to compile shader: %w", err)
	}
	return &PostProcessor{shader: shader}, nil
}

// Apply applies the CRT shader effect to the screen
func (pp *PostProcessor) Apply(screen *ebiten.Image) {
	w := screen.Bounds().Dx()
	h := screen.Bounds().Dy()

	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = screen
	opts.Uniforms = map[string]interface{}{
		// Kage uniforms can be added here if needed
	}

	img := ebiten.NewImage(w, h)
	img.DrawRectShader(w, h, pp.shader, opts)
	screen.DrawImage(img, nil)
}
