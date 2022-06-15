package enemy

import (
	"github.com/faiface/pixel"
	"math"
	"topdown/bullet"
	"topdown/globals"
	"topdown/player"
)

type Enemy struct {
	Rect     pixel.Rect
	Sheet    pixel.Picture
	Dir      float64
	Sprite   *pixel.Sprite
	Frame    pixel.Rect
	Anims    map[string][]pixel.Rect
	Counter  float64
	Rate     float64
	Speed    float64
	Life     int
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

		moveX := true
		moveY := true
		for angle := 0.0; angle < 2*math.Pi; angle += 0.10 {
			if globals.CheckWallCollision(e.Rect.Moved(pixel.V(p.Rect.Center().Sub(e.Rect.Center()).X, 0).Unit().Scaled(dt*e.Speed)).Center().X+(math.Cos(angle)*e.Rect.W()/3), e.Rect.Moved(pixel.V(p.Rect.Center().Sub(e.Rect.Center()).X, 0).Unit().Scaled(dt*e.Speed)).Center().Y+(math.Sin(angle)*e.Rect.H()/3), globals.ColisionImageSrc) == 1 {
				moveX = false
			}
			if globals.CheckWallCollision(e.Rect.Moved(pixel.V(0, p.Rect.Center().Sub(e.Rect.Center()).Y).Unit().Scaled(dt*e.Speed)).Center().X+(math.Cos(angle)*e.Rect.W()/3), e.Rect.Moved(pixel.V(0, p.Rect.Center().Sub(e.Rect.Center()).Y).Unit().Scaled(dt*e.Speed)).Center().Y+(math.Sin(angle)*e.Rect.H()/3), globals.ColisionImageSrc) == 1 {
				moveY = false
			}
		}

		if moveX {
			e.Rect = e.Rect.Moved(pixel.V(p.Rect.Center().Sub(e.Rect.Center()).Unit().Scaled(dt*e.Speed).X, 0))
		}
		if moveY {
			e.Rect = e.Rect.Moved(pixel.V(0, p.Rect.Center().Sub(e.Rect.Center()).Unit().Scaled(dt*e.Speed).Y))
		}
		e.Dir = p.Rect.Center().Sub(e.Rect.Center()).Angle()
	} else {
		e.Frame = e.Anims["death"][0]
	}
}

func (e *Enemy) CheckHit(b *bullet.Bullet) int {
	if e.Rect.Intersects(b.Rect) && e.Hittimer <= 0 {
		e.Life -= 1
		e.Hittimer = 0.5
		return 1
	}
	return 0
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
		Rotated(e.Rect.Center(), e.Dir+math.Pi/2),
	)
}
