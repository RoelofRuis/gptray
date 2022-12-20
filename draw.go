package main

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"time"
)

func Draw(w io.Writer, width, height int) error {
	rand.Seed(time.Now().UnixNano())

	// Image
	aspectRatio := float64(width) / float64(height)
	samplesPerPixel := 100
	maxDepth := 50

	// World
	world := World{
		Sphere2{Vector{0, 0, -1}, 0.5},
		Sphere2{Vector{0, -100.5, -1}, 100},
	}

	// Camera
	camera := NewCamera(aspectRatio)

	// Render
	_, err := fmt.Fprintf(w, "P3\n%d %d\n255\n", width, height)
	if err != nil {
		return err
	}

	for j := height - 1; j >= 0; j-- {
		progress := float64(height-j) / float64(height) * 100
		fmt.Printf("\rProgress: %d%%", int(progress))
		for i := 0; i < width; i++ {
			color := Color{}
			for s := 0; s < samplesPerPixel; s++ {
				x := (float64(i) + rand.Float64()) / float64(width-1)
				y := (float64(j) + rand.Float64()) / float64(height-1)

				ray := camera.GetRay(x, y)
				color = color.Add(RayColor(ray, world, maxDepth))
			}

			err := WriteColor(w, color, samplesPerPixel)
			if err != nil {
				return err
			}
		}
	}

	fmt.Print("\nDone\n")

	return nil
}

// RayColor returns the color for the given Ray
func RayColor(r Ray, world World, depth int) Color {
	// we've exceeded the ray bounce limit, no more light is gathered.
	if depth <= 0 {
		return Color{}
	}

	if rec, hasHit := world.Hit(r, 0.001, math.MaxFloat64); hasHit {
		target := rec.Position.Add(rec.Normal).Add(RandomUnitVector())
		return RayColor(
			Ray{rec.Position, target.Sub(rec.Position)},
			world,
			depth-1,
		).MulScalar(0.5)
	}

	unitDirection := r.Direction.Unit()
	t := 0.5 * (unitDirection.Y + 1.0)
	return Color{1, 1, 1}.MulScalar(1.0 - t).Add(Color{0.5, 0.7, 1.0}.MulScalar(t))
}
