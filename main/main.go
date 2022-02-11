package main

import (
	"errors"
	"flag"
	"fmt"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/theotheradamsmith/verbose-octo-eureka/logic"
)

func process_pixel(pixel color.Color) (uint8, error) {
	r, g, b, a := pixel.RGBA()
	if r == g && g == b && a == 65535 {
		return uint8(r * 255 / 65535), nil
	}
	return 0, errors.New("Invalid pixel coloration detected")
}

func decode(imagePath string) (string, error) {
	var imageStr strings.Builder

	f, ok := os.Open(imagePath)
	if ok != nil {
		return "", ok
	}

	defer f.Close()

	decodedImage, ok := png.Decode(f)
	if ok != nil {
		return "", ok
	}

	if decodedImage.Bounds().Dx() != 27 || decodedImage.Bounds().Dy() != 27 {
		msg := fmt.Sprintf(
			"Image provided is of invalid size %dx%d",
			decodedImage.Bounds().Dx(),
			decodedImage.Bounds().Dy(),
		)
		return "", errors.New(msg)
	}

	// upper-left:  (Bounds().Min.X, Bounds().Min.Y)
	// lower-right: (Bounds().Max.X-1, Bounds().Max.Y-1)
	for y := decodedImage.Bounds().Min.Y; y < decodedImage.Bounds().Max.Y; y += 3 {
		for x := decodedImage.Bounds().Min.X; x < decodedImage.Bounds().Max.X; x += 3 {
			pixel1 := decodedImage.At(x, y)
			pixel2 := decodedImage.At(x+1, y)
			pixel3 := decodedImage.At(x+2, y)
			if pixel1 == pixel2 && pixel2 == pixel3 {
				if p, ok := process_pixel(pixel1); ok == nil {
					if hexfmtint, ok := strconv.Atoi(fmt.Sprintf("%x", p)); ok != nil {
						return "", ok
					} else {
						letter := strconv.Itoa(hexfmtint % 10)
						imageStr.WriteString(letter)
					}
				} else {
					return "", ok
				}
			} else {
				msg := fmt.Sprintf("Mismatched pixels: {%v,%v,%v}", pixel1, pixel2, pixel3)
				return "", errors.New(msg)
			}
		}
	}

	return imageStr.String(), nil
}

func main() {
	fmt.Println("Hello, CTF!")
	pFlag := flag.String("path", "", "path of the image to decode")
	flag.Parse()
	if *pFlag != "" {
		object, ok := decode(*pFlag)
		if ok != nil {
			fmt.Println(ok)
			return
		}
		if _, ok := logic.Check(object); ok != nil {
			fmt.Println(ok)
		} else {
			fmt.Println("Congratulations! You have solved the puzzle!")
		}
	}
	// read from .config file
}
