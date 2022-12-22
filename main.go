package main

import (
	"os"
)

func main() {
	// Image
	image := NewImage(320, 240, 100, 50)

	// World
	world := World{}

	materialGround := Lambertian{Color{0.8, 0.8, 0.0}}
	materialCenter := Lambertian{Color{0.1, 0.2, 0.5}}
	materialLeft := Dielectric{1.5}
	materialRight := Metal{Color{0.8, 0.6, 0.2}, 0.0}

	world = append(world, Sphere2{Vector{0.0, -100.5, -1.0}, 100.0, materialGround})
	world = append(world, Sphere2{Vector{0.0, 0.0, -1.0}, 0.5, materialCenter})
	world = append(world, Sphere2{Vector{-1.0, 0.0, -1.0}, 0.5, materialLeft})
	world = append(world, Sphere2{Vector{1.0, 0.0, -1.0}, 0.5, materialRight})

	// Camera
	lookFrom := Vector{3, 3, 2}
	lookAt := Vector{0, 0, -1}
	camera := NewCamera(
		lookFrom,
		lookAt,
		Vector{0, 1, 0},
		20,
		image.AspectRatio,
		0.1,
		lookFrom.Sub(lookAt).Length(),
	)

	Draw(world, camera, image)

	file, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := image.Save(file); err != nil {
		panic(err)
	}
}
