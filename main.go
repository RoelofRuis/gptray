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

//func main() {
//	// define the scene
//	scene := &Scene{
//		AmbientColor:     &Color2{0, 0, 0},
//		AmbientIntensity: 0.1,
//		Camera2: Camera2{
//			Position: Vector{0, 0, 0},
//			LookAt:   Vector{0, 0, 1},
//			Up:       Vector{0, 1, 0},
//			Fov:      45,
//		},
//		Lights: []Light{
//			{
//				Position:  Vector{10, 10, 10},
//				Intensity: 1.0,
//				Color2:     &Color2{1, 1, 1},
//			},
//			{
//				Position:  Vector{-10, 10, 6},
//				Intensity: 0.4,
//				Color2:     &Color2{0, 0, 1},
//			},
//		},
//		Objects: []Object{
//			&Sphere{
//				Center: Vector{0, 0.2, 5},
//				Radius: 1,
//				MaterialProperties: Material{
//					Color2:            &Color2{1, 0, 0},
//					Specular:         0.01,
//					Diffuse:          0.8,
//					SpecularExponent: 4.0,
//					Reflective:       0.5,
//					Refractive:       0.0,
//					RefractionIndex:  1.0,
//				},
//			},
//			&Sphere{
//				Center: Vector{2.5, -0.5, 3},
//				Radius: 0.8,
//				MaterialProperties: Material{
//					Color2:            &Color2{0, 1, 0},
//					Specular:         0.01,
//					Diffuse:          0.5,
//					SpecularExponent: 4.0,
//					Reflective:       0.5,
//					Refractive:       0.0,
//					RefractionIndex:  33.0,
//				},
//			},
//			Plane{
//				Position:     Vector{0, -2, 0},
//				NormalVector: Vector{0, 1, 0},
//				MaterialProperties: Material{
//					Color2:            &Color2{0.5, 0.5, 0.5},
//					Specular:         0.01,
//					Diffuse:          1.0,
//					SpecularExponent: 4.0,
//					Reflective:       0.0,
//					Refractive:       0.0,
//					RefractionIndex:  1.0,
//				},
//			},
//		},
//	}
//
//	// generate the image
//	img := Raytrace(scene, 640, 480)
//
//	// save the image to disk
//	f, err := os.Create("image.png")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer f.Close()
//	png.Encode(f, img)
//}
