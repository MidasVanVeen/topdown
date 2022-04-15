package sheetloader

import (
	"io"
	"os"
	"image"
	"encoding/csv"
	"strconv"
	"github.com/faiface/pixel"
	_ "image/png"
	"github.com/pkg/errors"
)

func LoadSheet(path, dpath string, width float64) (sheet pixel.Picture, anims map[string][]pixel.Rect, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "error loading animation sheet")
		}
	}()
	file, err := os.Open(path) // open file
	if err != nil {
		return nil, nil, err
	}
	defer file.Close() // close after function 

	img, _, err := image.Decode(file) // decode as image
	if err != nil {
		return nil, nil, err
	}
	sheet = pixel.PictureDataFromImage(img)
	
	var frames []pixel.Rect
	for i := 0.0; i + width <= sheet.Bounds().Max.X; i += width {
		frames = append(frames, pixel.R(
			i,
			0,
			i+width,
			sheet.Bounds().H(),
		))
	}

	dfile, err := os.Open(dpath) // open file
	if err != nil {
		return nil, nil, err
	}
	defer dfile.Close() // close after function 

	anims = make(map[string][]pixel.Rect)

	csvfile := csv.NewReader(dfile)
	for {
		anim, err := csvfile.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		name := anim[0]
		start, _ := strconv.Atoi(anim[1])
		end, _ := strconv.Atoi(anim[2])

		anims[name] = frames[start : end+1]
	}

	return sheet, anims, nil
}
