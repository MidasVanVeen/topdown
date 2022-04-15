package bullet

import (
	"fmt"
	"github.com/faiface/pixel"
)

type Bullet struct {
	Vel pixel.Vec
	Rect pixel.Rect
	Life float64
	Sprite *pixel.Sprite
	Dir float64
	Sheet pixel.Picture
	Frame pixel.Rect
}

func (b *Bullet) Update(dt float64) {
	fmt.Println("bulletVel: ", b.Vel)
	b.Life -= dt
	b.Rect = b.Rect.Moved(b.Vel.Scaled(dt))
}

func (b *Bullet) Draw(t pixel.Target) {
	if b.Sprite == nil {
		b.Sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	b.Sprite.Set(b.Sheet, b.Frame)
	b.Sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			b.Rect.W()/b.Sprite.Frame().W(),
			b.Rect.H()/b.Sprite.Frame().H(),
		)).
		Moved(b.Rect.Center()).
		Rotated(b.Rect.Center(), b.Dir),
	)
}
