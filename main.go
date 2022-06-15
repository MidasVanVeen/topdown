package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math"
	"math/rand"
	"time"
	"topdown/bullet"
	"topdown/enemy"
	"topdown/globals"
	"topdown/player"
	"topdown/sheetloader"
	"os"
)

const (
	BULLETTIMEOUT     = 0.2
)

var (
	ENEMYSPAWNTIMEOUT = 5.0
)

func SpawnEnemy(x, y float64, enemysheet pixel.Picture, enemyanims map[string][]pixel.Rect) enemy.Enemy {
	return enemy.Enemy{
		Rect:  pixel.R(0, 0, 120*0.90, 96*0.90).Moved(pixel.V(x, y)),
		Sheet: enemysheet, Anims: enemyanims,
		Rate:  0.5,
		Speed: 80,
		Life:  3,
	}
}

func run() {
	rand.Seed(time.Now().UnixNano())

	cfg := pixelgl.WindowConfig{
		Title:  "topdown",
		Bounds: pixel.R(0, 0, 1024, 768),
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	globals.InitGlobals()

	// inladen van alle plaatjes en animaties dmv een sheetloader
	gameoversheet, gameoveranims, err := sheetloader.LoadSheet("assets/gameover/1.png", "assets/gameover/1.csv", 1000)
	if err != nil {
		panic(err)
	}
	gameoverRect := pixel.R(-1000, -1000, 1000, 1000)
	backgroundsheet, backgroundanims, err := sheetloader.LoadSheet("assets/levels/1.png", "assets/levels/1.csv", 1000)
	if err != nil {
		panic(err)
	}
	backgroundRect := pixel.R(-1000*1.5, -833*1.5, 1000*1.5, 833*1.5)
	overlaysheet, overlayanims, err := sheetloader.LoadSheet("assets/other/overlay.png", "assets/other/overlay.csv", 640)
	if err != nil {
		panic(err)
	}
	overlayRect := pixel.R(-640/2, -480/2, 640/2, 480/2)
	healthbarsheet, healthbaranims, err := sheetloader.LoadSheet("assets/healthbar/healthbar.png", "assets/healthbar/healthbar.csv", 540/5)
	if err != nil {
		panic(err)
	}
	healthbarRect := pixel.R(-540/5, -10, 540/5, 10)
	bulletsheet, bulletanims, err := sheetloader.LoadSheet("assets/bullet/bullet.png", "assets/bullet/bullet.csv", 87)
	if err != nil {
		panic(err)
	}
	enemysheet, enemyanims, err := sheetloader.LoadSheet("assets/zombie/zombie.png", "assets/zombie/zombie.csv", 100)
	if err != nil {
		panic(err)
	}
	playersheet, playeranims, err := sheetloader.LoadSheet("assets/player/playersheet.png", "assets/player/playersheet.csv", 100)
	if err != nil {
		panic(err)
	}

	// playerphysics object. hier gebeuren alle berekeningen voor het player object. denk aan beweging en controls
	playerphys := &player.PlayerPhysics{
		Vel:      pixel.V(0, 0),
		Rect:     pixel.R(0, 0, 120*0.90, 96*0.90),
		MaxSpeed: 220,
		State:    1,
		Life:     0,
		Score: 1,
	}
	// playeranimations object. hier gebeuren alle animatie gerelateerde dingen.
	playeranim := &player.PlayerAnimation{
		Sheet: playersheet,
		Anims: playeranims,
		Rate:  0.2,
		Dir:   1,
	}

	bullets := []bullet.Bullet{} // empty bullet array

	enemies := []enemy.Enemy{} // empty enemy array

	// TODO: make function for this
	for i := 0; i < 5; i++ {
		enemies = append(enemies, SpawnEnemy(float64(rand.Intn(1000)), float64(rand.Intn(1000)), enemysheet, enemyanims))
	}

	// canvas waar alle imd's op worden geprojecteerd. dit canvas kan verplaatst worden waardoor je de viewport verplaatst.
	canvas := pixelgl.NewCanvas(pixel.R(-640/2, -480/2, 640/2, 480/2))
	camPos := pixel.ZV // empty camera pos vector.

	// imd declarations, dit kan je zien als de render lagen van het spelletje.
	playerimd := imdraw.New(playersheet)
	playerimd.Precision = 32
	enemyimd := imdraw.New(enemysheet)
	enemyimd.Precision = 32
	bulletsimd := imdraw.New(bulletsheet)
	bulletsimd.Precision = 32
	backgroundimd := imdraw.New(backgroundsheet)
	backgroundimd.Precision = 32
	gameoverimd := imdraw.New(gameoversheet)
	gameoverimd.Precision = 32
	overlayimd := imdraw.New(overlaysheet)
	overlayimd.Precision = 32
	healthbarimd := imdraw.New(healthbarsheet)
	healthbarimd.Precision = 32

	// timer voor schieten. als deze 0 is mag de speler weer shieten.
	bullettimer := 0.0
	enemytimer := 0.0

	gameoverSprite := pixel.NewSprite(nil, pixel.Rect{})
	gameoverSprite.Set(gameoversheet, gameoveranims["main"][0])
	gameoverSprite.Draw(gameoverimd, pixel.IM.ScaledXY(pixel.ZV, pixel.V(
		gameoverRect.W()/gameoveranims["main"][0].W(),
		gameoverRect.H()/gameoveranims["main"][0].H(),
	)).Moved(gameoverRect.Center()))

	// eerste tijdspunt voor berekenen van deltatime
	t1 := time.Now()
	// laatste richting waar de speler naar keek. dit is alleen voor het gerekenen van de richting van de bullets
	lastdir := globals.Ctrl{1, 0, 0} // dit is van type Ctrl (control). zie topdown/types voor meer info
	for {
		// get deltatime
		t2 := time.Now()                               // tweede tijdspunt voor berekenen van deltatime.
		dt := t2.Sub(t1).Seconds()                     // deltatime in secondes. dit is de tijd sinds de laatste loop iteratie.
		fmt.Println("pos: ", playerphys.Rect.Center()) // DEBUG
		t1 = time.Now()                                // zet het eerste weer naar de huidige tijd

		fmt.Println(int(playerphys.Rect.Center().X/3+500), int(playerphys.Rect.Center().Y/3+833/2))

		// set cam pos
		camPos = pixel.Lerp(camPos, playerphys.Rect.Center(), 1-math.Pow(1.0/128, dt)) // lerp geeft een punt op een lijn tussen 2 punten
		cam := pixel.IM.Moved(camPos.Scaled(-1))
		canvas.SetMatrix(cam)

		// update
		// controls
		ctrl := globals.Ctrl{0, 0, 0}
		if win.Pressed(pixelgl.KeyLeft) {
			ctrl.X = -1
			lastdir = globals.Ctrl{-1, 0, 0}
		}
		if win.Pressed(pixelgl.KeyRight) {
			ctrl.X = 1
			lastdir = globals.Ctrl{1, 0, 0}
		}
		if win.Pressed(pixelgl.KeyDown) {
			ctrl.Y = -1
			lastdir = globals.Ctrl{0, -1, 0}
		}
		if win.Pressed(pixelgl.KeyUp) {
			ctrl.Y = 1
			lastdir = globals.Ctrl{0, 1, 0}
		}
		if win.Pressed(pixelgl.KeyS) {
			ctrl.S = 1
		}
		if win.Pressed(pixelgl.KeyQ) {
			os.Exit(0)
		}

		bullettimer -= dt // decrement van bullettimer met 1 per seconde
		if win.Pressed(pixelgl.KeySpace) && ctrl.S == 1 {
			if bullettimer <= 0 { // check of bullet
				bullets = append(bullets, bullet.Bullet{
					Vel:   pixel.V(lastdir.X*1200, lastdir.Y*1200),
					Rect:  pixel.R(-0.1*87/2, -0.1*184/2, 0.1*87/2, 0.1*184/2).Moved(playerphys.Rect.Center()).Moved(pixel.V(lastdir.X*40, lastdir.X*-20)).Moved(pixel.V(lastdir.Y*25, lastdir.Y*25)),
					Life:  0.5,
					Sheet: bulletsheet,
					Frame: bulletanims["main"][0],
					Dir:   playeranim.Dir + math.Pi,
				})
				bullettimer = BULLETTIMEOUT
			}
		}

		enemytimer -= dt
		if enemytimer <= 0 {
			enemies = append(enemies, SpawnEnemy(float64(rand.Intn(1000)), float64(rand.Intn(1000)), enemysheet, enemyanims))
			enemytimer = ENEMYSPAWNTIMEOUT
		}

		playerphys.Update(dt, ctrl)
		playeranim.Update(dt, playerphys, ctrl)

		for i := range enemies {
			if enemies[i].Life > 0 {
				playerphys.CheckEnemyHit(enemies[i].Rect)
			}
		}

		for i := range bullets {
			bullets[i].Update(dt)
		}
		for i := 0; i < len(bullets); i++ {
			if bullets[i].Life <= 0 {
				bullets = append(bullets[:i], bullets[i+1:]...)
			}
		}

		for i := 0; i < len(enemies); i++ {
			if enemies[i].Life <= 0 {
				enemies = append(enemies[:i], enemies[i+1:]...)
			}
		}

		for i := range enemies {
			enemies[i].Update(dt, playerphys)
			for _, b := range bullets {
				s := enemies[i].CheckHit(&b)
				if s == 1 {
					ENEMYSPAWNTIMEOUT *= 0.99
				}
			}
		}

		// clearr
		canvas.Clear(colornames.Black)
		playerimd.Clear()
		backgroundimd.Clear()
		bulletsimd.Clear()
		overlayimd.Clear()
		enemyimd.Clear()
		healthbarimd.Clear()

		// render
		backgroundSprite := pixel.NewSprite(nil, pixel.Rect{})
		backgroundSprite.Set(backgroundsheet, backgroundanims["main"][0])
		backgroundSprite.Draw(backgroundimd, pixel.IM.ScaledXY(pixel.ZV, pixel.V(
			backgroundRect.W()/backgroundanims["main"][0].W(),
			backgroundRect.H()/backgroundanims["main"][0].H(),
		)).Moved(backgroundRect.Center()))

		backgroundimd.Draw(canvas)

		for i := range bullets {
			bullets[i].Draw(bulletsimd)
		}

		bulletsimd.Draw(canvas)

		for i := range enemies {
			enemies[i].Draw(enemyimd)
		}

		if playerphys.Life >= 4 {
			fmt.Println("Game over")
			//gameoverimd.Draw(canvas)
			//canvas.SetMatrix(pixel.IM.Moved(gameoverRect.Center()))
			//canvas.Draw(win, pixel.IM.Moved(gameoverRect.Center()))
			playerphys = &player.PlayerPhysics{
				Vel:      pixel.V(0, 0),
				Rect:     pixel.R(0, 0, 120*0.90, 96*0.90),
				MaxSpeed: 220,
				State:    1,
				Life:     0,
			}
			// playeranimations object. hier gebeuren alle animatie gerelateerde dingen.
			playeranim = &player.PlayerAnimation{
				Sheet: playersheet,
				Anims: playeranims,
				Rate:  0.2,
				Dir:   1,
			}
			enemies = []enemy.Enemy{}
			bullets = []bullet.Bullet{}
			ENEMYSPAWNTIMEOUT = 5.0
		} else {
			enemyimd.Draw(canvas)

			playeranim.Draw(playerimd, playerphys)
			playerimd.Draw(canvas)

			overlaySprite := pixel.NewSprite(nil, pixel.Rect{})
			overlaySprite.Set(overlaysheet, overlayanims["main"][0])
			overlaySprite.Draw(overlayimd, pixel.IM.ScaledXY(pixel.ZV, pixel.V(
				overlayRect.W()/overlayanims["main"][0].W(),
				overlayRect.H()/overlayanims["main"][0].H(),
			)).Moved(camPos))

			overlayimd.Draw(canvas)

			healthbarSprite := pixel.NewSprite(nil, pixel.Rect{})
			healthbarSprite.Set(healthbarsheet, healthbaranims["main"][playerphys.Life])
			healthbarSprite.Draw(healthbarimd, pixel.IM.ScaledXY(pixel.ZV, pixel.V(
				healthbarRect.W()/healthbaranims["main"][playerphys.Life].W(),
				healthbarRect.H()/healthbaranims["main"][playerphys.Life].H(),
			)).Moved(camPos.Add(pixel.V(-200, 200))))

			healthbarimd.Draw(canvas)

			win.Clear(colornames.White)
			win.SetMatrix(pixel.IM.Scaled(pixel.ZV,
				math.Min(
					win.Bounds().W()/canvas.Bounds().W(),
					win.Bounds().H()/canvas.Bounds().H(),
				),
			).Moved(win.Bounds().Center()))
			canvas.Draw(win, pixel.IM.Moved(canvas.Bounds().Center()))
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
