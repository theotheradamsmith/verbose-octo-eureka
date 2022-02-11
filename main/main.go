package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/theotheradamsmith/verbose-octo-eureka/image"
	"github.com/theotheradamsmith/verbose-octo-eureka/logic"
)

func main() {
	fmt.Println("Hello, CTF!")
	pFlag := flag.String("path", "", "path of the image to decode")
	flag.Parse()
	if *pFlag != "" {
		f, ok := os.Open(*pFlag)
		if ok != nil {
			fmt.Println(ok)
		}
		defer f.Close()
		object, ok := image.Decode(f)
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
