package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func Draw(image *Image) {
	rand.Seed(time.Now().UnixNano())

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

	for j := 0; j < image.Height; j++ {
		progress := float64(j) / float64(image.Height) * 100
		fmt.Printf("\rProgress: %d%%", int(progress))
		for i := 0; i < image.Width; i++ {
			color := Color{}
			for s := 0; s < image.SamplesPerPixel; s++ {
				x := (float64(i) + rand.Float64()) / float64(image.Width-1)
				y := (float64(j) + rand.Float64()) / float64(image.Height-1)

				ray := camera.GetRay(x, y)
				color = color.Add(RayColor(ray, world, image.MaxDepth))
			}

			image.WriteColor(i, image.Height-j, color)
		}
	}
}

// RayColor returns the color for the given Ray
func RayColor(r Ray, world World, depth int) Color {
	// we've exceeded the ray bounce limit, no more light is gathered.
	if depth <= 0 {
		return Color{}
	}

	if rec, hasHit := world.Hit(r, 0.001, math.MaxFloat64); hasHit {
		isScattered, attenuation, scattered := rec.Material.Scatter(r, rec)

		if isScattered {
			return RayColor(scattered, world, depth-1).Mul(attenuation)
		}

		return Color{}
	}

	unitDirection := r.Direction.Unit()
	t := 0.5 * (unitDirection.Y + 1.0)
	return Color{1, 1, 1}.MulScalar(1.0 - t).Add(Color{0.5, 0.7, 1.0}.MulScalar(t))
}
