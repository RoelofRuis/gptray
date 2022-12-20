package main

import (
	"fmt"
	"io"
	"math"
)

type Color = Vector

func WriteColor(w io.Writer, c Color, samplesPerPixel int) error {
	scale := 1.0 / float64(samplesPerPixel)

	// Divide the color by the number of samples and gamma-correct for gamma-2.0
	r := math.Sqrt(c.X * scale)
	g := math.Sqrt(c.Y * scale)
	b := math.Sqrt(c.Z * scale)

	ir := int(256 * Clamp(r, 0.0, 0.999))
	ig := int(256 * Clamp(g, 0.0, 0.999))
	ib := int(256 * Clamp(b, 0.0, 0.999))

	// write the translated [0, 255] value of each color component.
	_, err := fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
	return err
}

func Clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
