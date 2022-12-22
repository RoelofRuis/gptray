package main

import (
	"os"
)

func main() {
	file, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	image := NewImage(320, 240, 100, 50)
	Draw(image)

	if err := image.Save(file); err != nil {
		panic(err)
	}
}
