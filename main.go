package main

import (
	"image/png"
	"log"
	"os"
)

func main() {
	// define the scene
	scene := &Scene{
		BackgroundColor: &Color{0, 0, 0},
		AmbientLight:    0.2,
		Camera: Camera{
			Position: Vector{0, 0, 0},
			LookAt:   Vector{0, 0, 1},
			Up:       Vector{0, 1, 0},
			Fov:      45,
		},
	}

	sphere := &Sphere{
		Center: Vector{0, 0, 5},
		Radius: 1,
		MaterialProperties: Material{
			Color:       &Color{1, 0, 0},
			Specular:    0.5,
			Shininess:   32,
			Reflective:  0.1,
			Transparent: 0,
			Refraction:  0,
		},
	}

	scene.Objects = append(scene.Objects, sphere)

	plane := Plane{
		Position:     Vector{0, -10, 0},
		NormalVector: Vector{0, 1, 0},
		MaterialProperties: Material{
			Color:       &Color{0.5, 0.5, 0.5},
			Specular:    0.1,
			Shininess:   3,
			Reflective:  0.5,
			Transparent: 0,
			Refraction:  0,
		},
	}

	scene.Objects = append(scene.Objects, plane)

	light := Light{
		Position:  Vector{10, 10, 10},
		Intensity: 1,
		Color:     &Color{1, 1, 1},
	}

	scene.Lights = append(scene.Lights, light)

	// generate the image
	img := Raytrace(scene, 640, 480, 50)

	// save the image to disk
	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
