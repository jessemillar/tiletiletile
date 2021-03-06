package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"
)

type boundaries struct {
	width  int
	height int
	center coordinates
	offset coordinates
}

type coordinates struct {
	x int
	y int
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: tiletiletile [input filename] [output filename] [size] [width (optional)] [height (optional)]")
		return
	}

	input := os.Args[1]
	output := os.Args[2]
	size := os.Args[3]

	if output[len(output)-3:] != "png" {
		fmt.Println("tiletiletile only supports creating PNG images")
		return
	}

	bounds := boundaries{}

	if size == "custom" {
		if len(os.Args) < 5 {
			fmt.Println("If \"custom\" is used, [width] and [height] must be supplied")
			return
		}

		width, err := strconv.Atoi(os.Args[4])
		if err != nil {
			log.Print(err)
			return
		}

		bounds.width = width

		height, err := strconv.Atoi(os.Args[5])
		if err != nil {
			log.Print(err)
			return
		}

		bounds.height = height
	} else if size == "iphone6" {
		bounds.width = 852
		bounds.height = 1608
	}

	bounds.center.x = bounds.width / 2
	bounds.center.y = bounds.height / 2

	rawImageFile, err := os.Open(input)
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

	// Add 1 to account for centering the pattern
	across := math.Ceil(float64(bounds.width)/float64(imageFile.Bounds().Dx())) + 1
	down := math.Ceil(float64(bounds.height)/float64(imageFile.Bounds().Dy())) + 1

	bounds.offset.x = imageFile.Bounds().Dx() - int(math.Mod(float64(bounds.center.x)-float64(imageFile.Bounds().Dx())/2, float64(imageFile.Bounds().Dx())))
	bounds.offset.y = imageFile.Bounds().Dy() - int(math.Mod(float64(bounds.center.y)-float64(imageFile.Bounds().Dy())/2, float64(imageFile.Bounds().Dy())))

	for rows := 0; rows < int(across); rows++ {
		for columns := 0; columns < int(down); columns++ {
			draw.Draw(canvas, image.Rect(imageFile.Bounds().Dx()*rows-bounds.offset.x, imageFile.Bounds().Dy()*columns-bounds.offset.y, bounds.width, bounds.height), imageFile, image.Point{0, 0}, draw.Src)
		}
	}

	exportedImage, err := os.Create(output)
	if err != nil {
		log.Print(err)
		return
	}

	defer exportedImage.Close()

	png.Encode(exportedImage, canvas)
}
