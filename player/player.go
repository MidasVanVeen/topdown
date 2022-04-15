package player

import (
	"fmt"
	"math"
	"github.com/faiface/pixel"
	"topdown/types"
)

const (
	drag = 400
	idle = 1
	running = 2
	runshooting = 3
	idleshooting = 4
)

type PlayerPhysics struct {
	Rect pixel.Rect
	Vel pixel.Vec
	MaxSpeed float64
	State int
	Life int
	Hittimer float64
}

type PlayerAnimation struct {
	Sheet pixel.Picture
	Anims map[string][]pixel.Rect
	Rate float64
	Counter float64
	Dir float64
	Frame pixel.Rect
	Sprite *pixel.Sprite
}

func (p *PlayerPhysics) Update(dt float64, ctrl types.Ctrl) {
	p.Hittimer -= dt
	if ctrl.X != 0 || ctrl.Y != 0 && ctrl.S == 0 {
		p.State = 2
	} 
	if ctrl.S != 0 {
		p.State = 4
		p.MaxSpeed = 110
	} else {
		p.MaxSpeed = 220
	}
	if (ctrl.X != 0 || ctrl.Y != 0) && ctrl.S != 0 {
		p.State = 3
	}
	ectrl := types.Ctrl{0,0,0}
	if ctrl == ectrl {
		p.State = 1
	}

	fmt.Println("state: ",p.State)
	fmt.Println("ctrl: ",ctrl)

	p.Vel.X += ctrl.X * dt * 800
	p.Vel.Y += ctrl.Y * dt * 800

	if p.Vel.X > p.MaxSpeed {
		p.Vel.X = p.MaxSpeed
	}
	if p.Vel.Y > p.MaxSpeed {
		p.Vel.Y = p.MaxSpeed
	}
	if p.Vel.X < -p.MaxSpeed {
		p.Vel.X = -p.MaxSpeed
	}
	if p.Vel.Y < -p.MaxSpeed {
		p.Vel.Y = -p.MaxSpeed
	}

	if p.Vel.X > drag * dt {
		p.Vel.X -= drag * dt
	}
	if p.Vel.Y > drag * dt {
		p.Vel.Y -= drag * dt
	}
	if p.Vel.X < -drag * dt {
		p.Vel.X += drag * dt
	}
	if p.Vel.Y < -drag * dt {
		p.Vel.Y += drag * dt
	}

	p.Rect = p.Rect.Moved(p.Vel.Scaled(dt))
}

func (p *PlayerPhysics) CheckEnemyHit(r pixel.Rect) {
	if p.Rect.Intersects(r) && p.Life < 4 && p.Hittimer <= 0 {
		p.Life += 1
		p.Hittimer = 1
	}
}

func (p *PlayerAnimation) Update(dt float64, phys *PlayerPhysics, ctrl types.Ctrl) {
	p.Counter += dt
	if ctrl.X == 1 {
		p.Dir = math.Pi/2
	}
	if ctrl.X == -1 {
		p.Dir = -math.Pi/2
	}
	if ctrl.Y == 1 {
		p.Dir = math.Pi
	}
	if ctrl.Y == -1 {
		p.Dir = math.Pi*2
	}

	if p.Counter >= 5 {
		p.Counter = 0
	}
	switch phys.State {
	case idle:
		p.Frame = p.Anims["idle"][0]
		p.Counter = 1
		p.Rate = 0.2
	case running:
		p.Rate = 0.2
		if p.Rate > 0.15 {
			p.Rate -= 0.02 * dt
		}
		i := int(math.Floor(p.Counter / p.Rate))
		p.Frame = p.Anims["running"][i%len(p.Anims["running"])]
	case runshooting:
		p.Rate = 0.4
		i := int(math.Floor(p.Counter / p.Rate))
		p.Frame = p.Anims["shooting"][i%len(p.Anims["shooting"])]
	case idleshooting:
		p.Frame = p.Anims["idleshooting"][0]
	}
}

func (p *PlayerAnimation) Draw(t pixel.Target, phys *PlayerPhysics) {
	if p.Sprite == nil {
		p.Sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	// draw the correct frame with the correct position and direction
	p.Sprite.Set(p.Sheet, p.Frame)
	p.Sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			phys.Rect.W()/p.Sprite.Frame().W(),
			phys.Rect.H()/p.Sprite.Frame().H(),
		)).
		Moved(phys.Rect.Center()).
		Rotated(phys.Rect.Center(), p.Dir),
	)
} 
