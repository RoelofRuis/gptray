package main

import (
	"os"
)

func main() {
	// Image
	image := NewImage(320, 240, 100, 50)

	// World
	checkerTexture := CheckerTexture{NewSolidColor(0.2, 0.3, 0.1), NewSolidColor(0.9, 0.9, 0.9)}

	materialGround := Lambertian{checkerTexture}
	materialCenter := Lambertian{NewSolidColor(0.1, 0.2, 0.5)}
	materialLeft := Dielectric{1.5}
	materialRight := Metal{Color{0.8, 0.6, 0.2}, 0.0}

	var hittables []Hittable
	hittables = append(hittables, Sphere2{Vector{0.0, -100.5, -1.0}, 100.0, materialGround})
	hittables = append(hittables, Sphere2{Vector{0.0, 0.0, -3.0}, 0.5, materialCenter})
	hittables = append(hittables, Sphere2{Vector{-1.0, 0.0, -1.0}, 0.5, materialLeft})
	hittables = append(hittables, Sphere2{Vector{0.0, 0.0, -5.0}, 0.5, materialRight})

	world := World{
		Background: Color{0.7, 0.8, 1.0},
		Hittables:  hittables,
	}

	// Camera
	lookFrom := Vector{-1, 1, 1}
	lookAt := Vector{-1, 0, -1}
	camera := NewCamera(
		lookFrom,
		lookAt,
		Vector{0, 1, 0},
		50,
		image.AspectRatio,
		0.01,
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
