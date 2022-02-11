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
)

//"github.com/theotheradamsmith/verbose-octo-eureka/logic"

func process_pixel(pixel color.Color) (uint8, error) {
	r, g, b, a := pixel.RGBA()
	if r == g && r == b && a == 65535 {
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

	//fmt.Println("Decoded Image: ", decodedImage)
	//fmt.Println("Bounds: ", decodedImage.Bounds())
	//fmt.Println("", decodedImage.ColorModel())
	//fmt.Printf("%T\n", decodedImage.Bounds())

	if decodedImage.Bounds().Dx() != 27 || decodedImage.Bounds().Dy() != 27 {
		msg := fmt.Sprintf(
			"Image provided is of invalid size %dx%d",
			decodedImage.Bounds().Dx(),
			decodedImage.Bounds().Dy(),
		)
		return "", errors.New(msg)
	}
	//fmt.Printf("decodedImage.Bounds().Dx: %d\n", decodedImage.Bounds().Dx())
	//fmt.Printf("decodedImage.Bounds().Dy: %d\n", decodedImage.Bounds().Dy())

	// upper-left:  (Bounds().Min.X, Bounds().Min.Y)
	// lower-right: (Bounds().Max.X-1, Bounds().Max.Y-1)
	for y := decodedImage.Bounds().Min.Y; y < decodedImage.Bounds().Max.Y; y += 3 {
		for x := decodedImage.Bounds().Min.X; x < decodedImage.Bounds().Max.X; x += 3 {
			if p, ok := process_pixel(decodedImage.At(x, y)); ok == nil {
				hexfmtstr := fmt.Sprintf("%x", p)
				hexfmtint, _ := strconv.Atoi(hexfmtstr)
				//fmt.Printf("%v", hexfmtint%10)
				letter := strconv.Itoa(hexfmtint % 10)
				imageStr.WriteString(letter)
			}

		}
		//fmt.Println()
	}

	return imageStr.String(), nil
}

func main() {
	fmt.Println("Hello, CTF!")
	pFlag := flag.String("path", "", "path of the image to decode")
	flag.Parse()
	_, ok := decode(*pFlag)
	if ok != nil {
		fmt.Println(ok)
		return
	}
}
