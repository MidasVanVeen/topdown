package globals

import (
	"image"
	"os"
)

type Ctrl struct {
	X float64
	Y float64
	S int
}

func CheckWallCollision(x,y float64,img image.Image) int {
	ximg := int(x/3+500)
	yimg := int(-y/3+833/2)
	r,_,_,_ := img.At(ximg,yimg).RGBA() 
	if r == 65535 {
		return 1
	}
	return 0
}

var ColisionImageSrc image.Image

func InitGlobals() {
	colisionImageFile, err := os.Open("assets/levels/1.walls.png")
	if err != nil {
		panic(err)
	}
	defer colisionImageFile.Close()
	ColisionImageSrc, _, err = image.Decode(colisionImageFile)
	if err != nil {
		panic(err)
	}
}
