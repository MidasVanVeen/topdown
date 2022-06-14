package gameover

import (
	"github.com/faiface/pixel"
)

struct Gameover {
	Rect pixel.Rect
	Sheet pixel.Picture
	Sprite pixel.Sprite
}

func (g *Gameover) Draw(t pixel.Target) {
	if g.Sprite == nil {
		g.Sprite = pixel.NewSprite()
	}
	g.Sprite
}