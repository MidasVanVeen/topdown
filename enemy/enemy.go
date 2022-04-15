package enemy

import (
	"math"
	"github.com/faiface/pixel"
	"topdown/player"
	"topdown/bullet"
)

type Enemy struct {
	Rect pixel.Rect
	Sheet pixel.Picture
	Dir float64
	Sprite *pixel.Sprite
	Frame pixel.Rect
	Anims map[string][]pixel.Rect
	Counter float64
	Rate float64
	Speed float64
	Life int
	Hittimer float64
}

func (e *Enemy) Update(dt float64, p *player.PlayerPhysics) {
	e.Hittimer -= dt
	if e.Life > 0 {
		e.Counter += dt
		if e.Counter >= 2 {
			e.Counter = 0
		}
		i := int(math.Floor(e.Counter / e.Rate))
		e.Frame = e.Anims["walk"][i%len(e.Anims["walk"])]
		e.Rect = e.Rect.Moved(p.Rect.Center().Sub(e.Rect.Center()).Unit().Scaled(dt * e.Speed)) 
		e.Dir = p.Rect.Center().Sub(e.Rect.Center()).Angle()
	} else {
		e.Frame = e.Anims["death"][0]
	}
}

func (e *Enemy) CheckHit(b *bullet.Bullet) {
	if e.Rect.Intersects(b.Rect) && e.Hittimer <= 0 {
		e.Life -= 1
		e.Hittimer = 0.5
	}
}

func (e *Enemy) Draw(t pixel.Target) {
	if e.Sprite == nil {
		e.Sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	// draw the correct frame with the correct position and direction
	e.Sprite.Set(e.Sheet, e.Frame)
	e.Sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			e.Rect.W()/e.Sprite.Frame().W(),
			e.Rect.H()/e.Sprite.Frame().H(),
		)).
		Moved(e.Rect.Center()).
		Rotated(e.Rect.Center(), e.Dir + math.Pi/2),
	)
}
