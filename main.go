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

type Pixel struct {
	red float64
	green float64
	blue float64
	alpha float64
}

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
	bounds := photo.Bounds()

	maxX := bounds.Max.X
	maxY := bounds.Max.Y

	if maxX <= 1000 && maxY <= 1000 {

		var sourceImage [1000][1000]Pixel

		var destinationImage [1000][1000]Pixel

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				rOriginal, gOriginal, bOriginal, aOriginal := colorToDoubles(photo.At(imageX, imageY))
				sourceImage[imageX][imageY].red = rOriginal
				sourceImage[imageX][imageY].green = gOriginal
				sourceImage[imageX][imageY].blue = bOriginal
				sourceImage[imageX][imageY].alpha = aOriginal
			}
		}

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				pixel := sourceImage[imageX][imageY]
				originalPixel := sourceImage[imageX][imageY]

				total := 1.0

				for offsetX := -5; offsetX < 6; offsetX++ {
					for offsetY := -5; offsetY < 6; offsetY++ {

						x := imageX + offsetX
						y := imageY + offsetY

						if x >= 0 && x < maxX && y >= 0 && y < maxY {

							newPixel := sourceImage[x][y]

							rPart := originalPixel.red - newPixel.red
							if rPart < 0 {
								rPart = -rPart
							}

							gPart := originalPixel.green - newPixel.green
							if gPart < 0 {
								gPart = -gPart
							}

							bPart := originalPixel.blue - newPixel.blue
							if bPart < 0 {
								bPart = -bPart
							}

							aPart := originalPixel.alpha - newPixel.alpha
							if aPart < 0 {
								aPart = -aPart
							}

							if (rPart + gPart + bPart + aPart) < 60 {
								distance := distanceMatrix[offsetX+5][offsetY+5]

								total += distance

								pixel.red += newPixel.red * distance
								pixel.green += newPixel.green * distance
								pixel.blue += newPixel.blue * distance
								pixel.alpha += newPixel.alpha * distance
							}
						}
					}
				}

				pixel.red /= total
				pixel.green /= total
				pixel.blue /= total
				pixel.alpha /= total

				destinationImage[imageX][imageY] = pixel
			}
		}

		drawableImage := image.NewRGBA(photo.Bounds())

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				pixel := destinationImage[imageX][imageY]

				drawableImage.Set(imageX, imageY,
					color.RGBA{uint8(pixel.red), uint8(pixel.green), uint8(pixel.blue), uint8(pixel.alpha)})
			}
		}

		return drawableImage
	} else if maxX <= 2000 && maxY <= 2000 {

		var sourceImage [2000][2000]Pixel

		var destinationImage [2000][2000]Pixel

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				rOriginal, gOriginal, bOriginal, aOriginal := colorToDoubles(photo.At(imageX, imageY))
				sourceImage[imageX][imageY].red = rOriginal
				sourceImage[imageX][imageY].green = gOriginal
				sourceImage[imageX][imageY].blue = bOriginal
				sourceImage[imageX][imageY].alpha = aOriginal
			}
		}

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				pixel := sourceImage[imageX][imageY]
				originalPixel := sourceImage[imageX][imageY]

				total := 1.0

				for offsetX := -5; offsetX < 6; offsetX++ {
					for offsetY := -5; offsetY < 6; offsetY++ {

						x := imageX + offsetX
						y := imageY + offsetY

						if x >= 0 && x < maxX && y >= 0 && y < maxY {

							newPixel := sourceImage[x][y]

							rPart := originalPixel.red - newPixel.red
							if rPart < 0 {
								rPart = -rPart
							}

							gPart := originalPixel.green - newPixel.green
							if gPart < 0 {
								gPart = -gPart
							}

							bPart := originalPixel.blue - newPixel.blue
							if bPart < 0 {
								bPart = -bPart
							}

							aPart := originalPixel.alpha - newPixel.alpha
							if aPart < 0 {
								aPart = -aPart
							}

							if (rPart + gPart + bPart + aPart) < 60 {
								distance := distanceMatrix[offsetX+5][offsetY+5]

								total += distance

								pixel.red += newPixel.red * distance
								pixel.green += newPixel.green * distance
								pixel.blue += newPixel.blue * distance
								pixel.alpha += newPixel.alpha * distance
							}
						}
					}
				}

				pixel.red /= total
				pixel.green /= total
				pixel.blue /= total
				pixel.alpha /= total

				destinationImage[imageX][imageY] = pixel
			}
		}

		drawableImage := image.NewRGBA(photo.Bounds())

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				pixel := destinationImage[imageX][imageY]

				drawableImage.Set(imageX, imageY,
					color.RGBA{uint8(pixel.red), uint8(pixel.green), uint8(pixel.blue), uint8(pixel.alpha)})
			}
		}

		return drawableImage
	} else if maxX <= 4000 && maxY <= 4000 {

		var sourceImage [4000][4000]Pixel

		var destinationImage [4000][4000]Pixel

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				rOriginal, gOriginal, bOriginal, aOriginal := colorToDoubles(photo.At(imageX, imageY))
				sourceImage[imageX][imageY].red = rOriginal
				sourceImage[imageX][imageY].green = gOriginal
				sourceImage[imageX][imageY].blue = bOriginal
				sourceImage[imageX][imageY].alpha = aOriginal
			}
		}

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				pixel := sourceImage[imageX][imageY]
				originalPixel := sourceImage[imageX][imageY]

				total := 1.0

				for offsetX := -5; offsetX < 6; offsetX++ {
					for offsetY := -5; offsetY < 6; offsetY++ {

						x := imageX + offsetX
						y := imageY + offsetY

						if x >= 0 && x < maxX && y >= 0 && y < maxY {

							newPixel := sourceImage[x][y]

							rPart := originalPixel.red - newPixel.red
							if rPart < 0 {
								rPart = -rPart
							}

							gPart := originalPixel.green - newPixel.green
							if gPart < 0 {
								gPart = -gPart
							}

							bPart := originalPixel.blue - newPixel.blue
							if bPart < 0 {
								bPart = -bPart
							}

							aPart := originalPixel.alpha - newPixel.alpha
							if aPart < 0 {
								aPart = -aPart
							}

							if (rPart + gPart + bPart + aPart) < 60 {
								distance := distanceMatrix[offsetX+5][offsetY+5]

								total += distance

								pixel.red += newPixel.red * distance
								pixel.green += newPixel.green * distance
								pixel.blue += newPixel.blue * distance
								pixel.alpha += newPixel.alpha * distance
							}
						}
					}
				}

				pixel.red /= total
				pixel.green /= total
				pixel.blue /= total
				pixel.alpha /= total

				destinationImage[imageX][imageY] = pixel
			}
		}

		drawableImage := image.NewRGBA(photo.Bounds())

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				pixel := destinationImage[imageX][imageY]

				drawableImage.Set(imageX, imageY,
					color.RGBA{uint8(pixel.red), uint8(pixel.green), uint8(pixel.blue), uint8(pixel.alpha)})
			}
		}

		return drawableImage
	} else if maxX <= 10000 && maxY <= 10000 {

		var sourceImage [10000][10000]Pixel

		var destinationImage [10000][10000]Pixel

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				rOriginal, gOriginal, bOriginal, aOriginal := colorToDoubles(photo.At(imageX, imageY))
				sourceImage[imageX][imageY].red = rOriginal
				sourceImage[imageX][imageY].green = gOriginal
				sourceImage[imageX][imageY].blue = bOriginal
				sourceImage[imageX][imageY].alpha = aOriginal
			}
		}

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				pixel := sourceImage[imageX][imageY]
				originalPixel := sourceImage[imageX][imageY]

				total := 1.0

				for offsetX := -5; offsetX < 6; offsetX++ {
					for offsetY := -5; offsetY < 6; offsetY++ {

						x := imageX + offsetX
						y := imageY + offsetY

						if x >= 0 && x < maxX && y >= 0 && y < maxY {

							newPixel := sourceImage[x][y]

							rPart := originalPixel.red - newPixel.red
							if rPart < 0 {
								rPart = -rPart
							}

							gPart := originalPixel.green - newPixel.green
							if gPart < 0 {
								gPart = -gPart
							}

							bPart := originalPixel.blue - newPixel.blue
							if bPart < 0 {
								bPart = -bPart
							}

							aPart := originalPixel.alpha - newPixel.alpha
							if aPart < 0 {
								aPart = -aPart
							}

							if (rPart + gPart + bPart + aPart) < 60 {
								distance := distanceMatrix[offsetX+5][offsetY+5]

								total += distance

								pixel.red += newPixel.red * distance
								pixel.green += newPixel.green * distance
								pixel.blue += newPixel.blue * distance
								pixel.alpha += newPixel.alpha * distance
							}
						}
					}
				}

				pixel.red /= total
				pixel.green /= total
				pixel.blue /= total
				pixel.alpha /= total

				destinationImage[imageX][imageY] = pixel
			}
		}

		drawableImage := image.NewRGBA(photo.Bounds())

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				pixel := destinationImage[imageX][imageY]

				drawableImage.Set(imageX, imageY,
					color.RGBA{uint8(pixel.red), uint8(pixel.green), uint8(pixel.blue), uint8(pixel.alpha)})
			}
		}

		return drawableImage
	} else if maxX <= 20000 && maxY <= 20000 {

		var sourceImage [20000][20000]Pixel

		var destinationImage [20000][20000]Pixel

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				rOriginal, gOriginal, bOriginal, aOriginal := colorToDoubles(photo.At(imageX, imageY))
				sourceImage[imageX][imageY].red = rOriginal
				sourceImage[imageX][imageY].green = gOriginal
				sourceImage[imageX][imageY].blue = bOriginal
				sourceImage[imageX][imageY].alpha = aOriginal
			}
		}

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				pixel := sourceImage[imageX][imageY]
				originalPixel := sourceImage[imageX][imageY]

				total := 1.0

				for offsetX := -5; offsetX < 6; offsetX++ {
					for offsetY := -5; offsetY < 6; offsetY++ {

						x := imageX + offsetX
						y := imageY + offsetY

						if x >= 0 && x < maxX && y >= 0 && y < maxY {

							newPixel := sourceImage[x][y]

							rPart := originalPixel.red - newPixel.red
							if rPart < 0 {
								rPart = -rPart
							}

							gPart := originalPixel.green - newPixel.green
							if gPart < 0 {
								gPart = -gPart
							}

							bPart := originalPixel.blue - newPixel.blue
							if bPart < 0 {
								bPart = -bPart
							}

							aPart := originalPixel.alpha - newPixel.alpha
							if aPart < 0 {
								aPart = -aPart
							}

							if (rPart + gPart + bPart + aPart) < 60 {
								distance := distanceMatrix[offsetX+5][offsetY+5]

								total += distance

								pixel.red += newPixel.red * distance
								pixel.green += newPixel.green * distance
								pixel.blue += newPixel.blue * distance
								pixel.alpha += newPixel.alpha * distance
							}
						}
					}
				}

				pixel.red /= total
				pixel.green /= total
				pixel.blue /= total
				pixel.alpha /= total

				destinationImage[imageX][imageY] = pixel
			}
		}

		drawableImage := image.NewRGBA(photo.Bounds())

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				pixel := destinationImage[imageX][imageY]

				drawableImage.Set(imageX, imageY,
					color.RGBA{uint8(pixel.red), uint8(pixel.green), uint8(pixel.blue), uint8(pixel.alpha)})
			}
		}

		return drawableImage
	} else {

		//This image is too damn big

		var sourceImage [100000][100000]Pixel

		var destinationImage [100000][100000]Pixel

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				rOriginal, gOriginal, bOriginal, aOriginal := colorToDoubles(photo.At(imageX, imageY))
				sourceImage[imageX][imageY].red = rOriginal
				sourceImage[imageX][imageY].green = gOriginal
				sourceImage[imageX][imageY].blue = bOriginal
				sourceImage[imageX][imageY].alpha = aOriginal
			}
		}

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				pixel := sourceImage[imageX][imageY]
				originalPixel := sourceImage[imageX][imageY]

				total := 1.0

				for offsetX := -5; offsetX < 6; offsetX++ {
					for offsetY := -5; offsetY < 6; offsetY++ {

						x := imageX + offsetX
						y := imageY + offsetY

						if x >= 0 && x < maxX && y >= 0 && y < maxY {

							newPixel := sourceImage[x][y]

							rPart := originalPixel.red - newPixel.red
							if rPart < 0 {
								rPart = -rPart
							}

							gPart := originalPixel.green - newPixel.green
							if gPart < 0 {
								gPart = -gPart
							}

							bPart := originalPixel.blue - newPixel.blue
							if bPart < 0 {
								bPart = -bPart
							}

							aPart := originalPixel.alpha - newPixel.alpha
							if aPart < 0 {
								aPart = -aPart
							}

							if (rPart + gPart + bPart + aPart) < 60 {
								distance := distanceMatrix[offsetX+5][offsetY+5]

								total += distance

								pixel.red += newPixel.red * distance
								pixel.green += newPixel.green * distance
								pixel.blue += newPixel.blue * distance
								pixel.alpha += newPixel.alpha * distance
							}
						}
					}
				}

				pixel.red /= total
				pixel.green /= total
				pixel.blue /= total
				pixel.alpha /= total

				destinationImage[imageX][imageY] = pixel
			}
		}

		drawableImage := image.NewRGBA(photo.Bounds())

		for imageX := 0; imageX < maxX; imageX++ {
			for imageY := 0; imageY < maxY; imageY++ {
				pixel := destinationImage[imageX][imageY]

				drawableImage.Set(imageX, imageY,
					color.RGBA{uint8(pixel.red), uint8(pixel.green), uint8(pixel.blue), uint8(pixel.alpha)})
			}
		}

		return drawableImage
	}
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
