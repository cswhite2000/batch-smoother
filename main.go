package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var distanceMatrix [11][11]float64
var length int
var finished = 0
var files chan string

func main() {

	paths := getPhotoFiles()

	length = len(paths)

	files = make(chan string, length)

	for _, path := range paths {
		files <- path
	}

	if length == 1 {
		fmt.Printf("Starting conversion of 1 image\n")
	} else {
		fmt.Printf("Starting conversion of %d images\n", length)
	}

	initDistance()

	var wg sync.WaitGroup

	os.Mkdir("output", 0770)

	numRoutines := length
	numCPU := runtime.NumCPU()

	if length > numCPU {
		numRoutines = numCPU
	}

	fmt.Printf("%d %d", numRoutines, numCPU)

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go handleFiles(&wg)
	}
	wg.Wait()
}

func initDistance() {
	for x := -5; x < 6; x++ {
		for y := -5; y < 6; y++ {
			distanceMatrix[x+5][y+5] = math.Pow(2, -(math.Sqrt(float64((x * x) + (y * y)))))
		}
	}
}

func smoothPhoto(photo image.Image) image.Image {

	drawableImage := image.NewRGBA(photo.Bounds())

	bounds := photo.Bounds()

	for imageX := 0; imageX < bounds.Max.X; imageX++ {
		for imageY := 0; imageY < bounds.Max.Y; imageY++ {

			r, g, b, a := colorToDoubles(photo.At(imageX, imageY))
			rOriginal, gOriginal, bOriginal, aOriginal := colorToDoubles(photo.At(imageX, imageY))

			total := 1.0

			for offsetX := -5; offsetX < 6; offsetX++ {
				for offsetY := -5; offsetY < 6; offsetY++ {

					x := imageX + offsetX
					y := imageY + offsetY

					if x >= 0 && x < bounds.Max.X && y >= 0 && y < bounds.Max.Y {

						rNew, gNew, bNew, aNew := colorToDoubles(photo.At(x, y))

						if (math.Abs(rOriginal-rNew)+math.Abs(gOriginal-gNew)+math.Abs(bOriginal-bNew)+math.Abs(aOriginal-aNew))/4.0 < 15 {

							distance := distanceMatrix[offsetX+5][offsetY+5]

							total += distance

							r += rNew * distance
							g += gNew * distance
							b += bNew * distance
							a += aNew * distance
						}
					}
				}
			}

			drawableImage.Set(imageX, imageY,
				color.RGBA{uint8(r / total), uint8(g / total), uint8(b / total), uint8(a / total)})

		}
	}

	return drawableImage
}

func colorToDoubles(color color.Color) (r float64, g float64, b float64, a float64) {
	ri, gi, bi, ai := color.RGBA()

	return float64(ri / 256.0), float64(gi / 256.0), float64(bi / 256.0), float64(ai / 256.0)

}

func handleFiles(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		exit := false
		select {
		case path := <-files:
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}

			photo, err := jpeg.Decode(file)
			if err != nil {
				log.Fatal(err)
			}

			finalImage := smoothPhoto(photo)

			outputFile, err := os.Create("output/" + path)
			if err != nil {
				log.Fatal(err)
			}

			jpeg.Encode(outputFile, finalImage, &jpeg.Options{Quality: 99})

			finished++
			fmt.Printf("Finished photo %d of %d\n", finished, length)
		default:
			exit = true
			break
		}
		if exit {
			break
		}
	}
}

func getPhotoFiles() []string {
	files, err := filepath.Glob("*.jpg")
	if err != nil {
		log.Fatal(err)
	}

	return files
}
