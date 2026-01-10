package rendering

import (
	"github.com/pbnjay/pixfont"
	"snakehem/pxterm16"
	"snakehem/pxterm24"
	"sync"
)

// FontManager manages font instances without globals
type FontManager struct {
	font16 *pixfont.PixFont
	font24 *pixfont.PixFont
	once16 sync.Once
	once24 sync.Once
}

// NewFontManager creates a new font manager
func NewFontManager() *FontManager {
	return &FontManager{}
}

// GetFont16 returns the 16px font (lazy initialization)
func (fm *FontManager) GetFont16() *pixfont.PixFont {
	fm.once16.Do(func() {
		fm.font16 = pxterm16.CreateFont()
	})
	return fm.font16
}

// GetFont24 returns the 24px font (lazy initialization)
func (fm *FontManager) GetFont24() *pixfont.PixFont {
	fm.once24.Do(func() {
		fm.font24 = pxterm24.CreateFont()
	})
	return fm.font24
}
