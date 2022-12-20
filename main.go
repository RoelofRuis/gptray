package main

import (
	"os"
)

func main() {
	file, err := os.Create("image.ppm")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = Draw(file, 320, 240)
	if err != nil {
		panic(err)
	}
}
