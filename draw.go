package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func Draw(world World, camera Camera, image Image) {
	rand.Seed(time.Now().UnixNano())

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
	// We've exceeded the ray bounce limit, no more light is gathered.
	if depth <= 0 {
		return Color{}
	}

	// If the ray hits nothing, return the background color.
	rec, hasHit := world.Hit(r, 0.001, math.MaxFloat64)
	if !hasHit {
		return world.Background
	}

	emitted := rec.Material.Emitted(rec.P, rec.U, rec.V)

	isScattered, attenuation, scattered := rec.Material.Scatter(r, rec)
	if !isScattered {
		return emitted
	}

	return RayColor(scattered, world, depth-1).Mul(attenuation).Add(emitted)
}
