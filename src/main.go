package main

import (
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func handleUploadPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	const maxUploadSize = 1024 * 1024
	r.ParseMultipartForm(maxUploadSize)
	var err string = ""

	file, _, ok := r.FormFile("file")
	if ok != nil {
		fmt.Fprintf(w, "<h1>File Error</h1><p>%s</p>", ok)
		err = ok.Error()
	} else {
		defer file.Close()
		if object, ok := Decode(file); ok != nil {
			err = ok.Error()
		} else {
			//fmt.Fprintf(w, object)
			if _, ok = Check(object); ok != nil {
				err = ok.Error()
			} else {
				fmt.Fprintf(w, "<h1>Congratulations!</h1><p>You have solved the puzzle!</p>")
			}
		}
	}
	fmt.Fprintf(w, err)
	//render(w, r, homepageTpl, "homepage_view", data)
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w\n", err))
	}

	port := fmt.Sprintf(":%s", viper.GetString("port"))
	if port == ":" {
		log.Print("Error in configuration file: 'port' not found. Defaulting to 8000.")
		port = ":8000"
	}

	router := mux.NewRouter()
	router.HandleFunc("/", handleUploadPost).Methods(http.MethodPost)
	http.Handle("/", router)
	log.Print("Listening on port ", port)
	log.Fatal(http.ListenAndServe(port, router))
}

const digits string = "123456789"
const rows string = "ABCDEFGHI"
const cols string = digits

var squares = cross(rows, cols)
var unitList = createUnitList(cols, rows)
var units = createUnits(squares, unitList)

//var peers = createPeers(units)

func Hello(name string) (string, error) {
	if name == "" {
		return name, errors.New("empty name")
	}

	message := fmt.Sprintf("Hello, %v! Thanks for stopping by to gander at the logic!", name)
	return message, nil
}

func cross(a, b string) []string {
	var ret []string
	for _, av := range a {
		for _, bv := range b {
			ret = append(ret, string(av)+string(bv))
		}
	}
	return ret
}

func createUnitList(cols, rows string) [][]string {
	ret := make([][]string, len(rows)*3)
	i := 0
	for _, c := range cols {
		ret[i] = cross(rows, string(c))
		i++
	}
	for _, r := range rows {
		ret[i] = cross(string(r), cols)
		i++
	}
	rs := []string{"ABC", "DEF", "GHI"}
	cs := []string{"123", "456", "789"}
	for _, r := range rs {
		for _, c := range cs {
			ret[i] = cross(r, string(c))
			i++
		}
	}
	return ret
}

func createUnits(squares []string, unitList [][]string) map[string][][]string {
	units := make(map[string][][]string, len(squares))
	for _, s := range squares {
		unit := make([][]string, 3)
		i := 0
		for _, u := range unitList {
			for _, su := range u {
				if s == su {
					unit[i] = u
					i++
					break
				}
			}
		}
		units[s] = unit
	}
	return units
}

func GridValues(grid string) (map[string]string, error) {
	// Convert grid into a dict of {square: char} with '0' or '.' for empties
	gridValues := make(map[string]string, len(squares))
	validChars := make([]string, 0, len(grid))

	for _, c := range grid {
		if strings.Contains(digits, string(c)) || strings.Contains(".0", string(c)) {
			validChars = append(validChars, string(c))
		}
	}

	if len(validChars) != 81 {
		return gridValues, errors.New("Invalid input grid")
	}

	for i, s := range squares {
		gridValues[s] = string(validChars[i])
	}

	return gridValues, nil
}

func Verify(values map[string]string) (bool, error) {
	unitSolved := func(unit []string) bool {
		digitsSet := make(map[string]bool, len(digits))
		for _, r := range digits {
			digitsSet[string(r)] = true
		}
		for _, s := range unit {
			key := string(values[s])
			if _, ok := digitsSet[key]; ok {
				delete(digitsSet, key)
			} else {
				return false
			}
		}
		return len(digitsSet) == 0
	}
	for _, unit := range unitList {
		if unitSolved(unit) != true {
			msg := fmt.Sprintf("Error in unit: [%s:%s]", unit[0], unit[8])
			return false, errors.New(msg)
		}
	}
	return true, nil
}

func Check(userInput string) (bool, error) {
	if gv, ok := GridValues(userInput); ok != nil {
		return false, ok
	} else {
		if _, ok := Verify(gv); ok != nil {
			return false, ok
		}
	}
	return true, nil
}

func process_pixel(pixel color.Color) (uint8, error) {
	r, g, b, a := pixel.RGBA()
	if r == g && g == b && a == 65535 {
		return uint8(r * 255 / 65535), nil
	}
	return 0, errors.New("Invalid pixel coloration detected")
}

func Decode(file io.Reader) (string, error) {
	var imageStr strings.Builder

	decodedImage, ok := png.Decode(file)
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
