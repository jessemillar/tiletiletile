package main

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
)

type boundaries struct {
	width  int
	height int
}

func main() {
	bounds := boundaries{}
	bounds.width = 800
	bounds.height = 600

	rawImageFile, err := os.Open("test.png")
	if err != nil {
		log.Print(err)
		return
	}

	defer rawImageFile.Close()

	imageFile, _, err := image.Decode(rawImageFile)
	if err != nil {
		log.Print(err)
		return
	}

	canvas := image.NewRGBA(image.Rect(0, 0, bounds.width, bounds.height))

	across := math.Ceil(float64(bounds.width) / float64(imageFile.Bounds().Dx()))
	down := math.Ceil(float64(bounds.height) / float64(imageFile.Bounds().Dy()))

	for rows := 0; rows < int(across); rows++ {
		for columns := 0; columns < int(down); columns++ {
			draw.Draw(canvas, image.Rect(imageFile.Bounds().Dx()*rows, imageFile.Bounds().Dy()*columns, bounds.width, bounds.height), imageFile, image.Point{0, 0}, draw.Src)
		}
	}

	exportedImage, err := os.Create("new.png")
	if err != nil {
		log.Print(err)
		return
	}

	defer exportedImage.Close()

	png.Encode(exportedImage, canvas)
}
