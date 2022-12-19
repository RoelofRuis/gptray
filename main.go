package main

import (
	"image/png"
	"log"
	"os"
)

func main() {
	// define the scene
	scene := &Scene{
		AmbientColor:     &Color{1, 1, 1},
		AmbientIntensity: 0.1,
		Camera: Camera{
			Position: Vector{0, 0, 0},
			LookAt:   Vector{0, 0, 1},
			Up:       Vector{0, 1, 0},
			Fov:      45,
		},
		Lights: []Light{
			{
				Position:  Vector{10, 10, 10},
				Intensity: 1,
				Color:     &Color{1, 1, 1},
			},
		},
		Objects: []Object{
			&Sphere{
				Center: Vector{0, 0, 5},
				Radius: 1,
				MaterialProperties: Material{
					Color:             &Color{1, 0, 0},
					SpecularIntensity: 1.0,
					DiffuseIntensity:  1.0,
					Shininess:         1.0,
					Reflectivity:      0.0,
					Transparency:      0.0,
					RefractionIndex:   1.0,
				},
			},
			&Sphere{
				Center: Vector{2.5, 0, 3},
				Radius: 0.8,
				MaterialProperties: Material{
					Color:             &Color{0, 1, 0},
					SpecularIntensity: 1.0,
					DiffuseIntensity:  1.0,
					Shininess:         1.0,
					Reflectivity:      0.0,
					Transparency:      0.0,
					RefractionIndex:   1.0,
				},
			},
			Plane{
				Position:     Vector{0, -10, 0},
				NormalVector: Vector{0, 1, 0},
				MaterialProperties: Material{
					Color:             &Color{0.5, 0.5, 0.5},
					SpecularIntensity: 1.0,
					DiffuseIntensity:  1.0,
					Shininess:         1.0,
					Reflectivity:      0.0,
					Transparency:      0.0,
					RefractionIndex:   1.0,
				},
			},
		},
	}

	// generate the image
	img := Raytrace(scene, 640, 480)

	// save the image to disk
	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
