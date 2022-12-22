package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
)

type Color = Vector

type Image struct {
	Width           int
	Height          int
	SamplesPerPixel int
	MaxDepth        int
	AspectRatio     float64
	Scale           float64
	image           *image.RGBA
}

func NewImage(width, height, samplesPerPixel, maxDepth int) Image {
	aspectRatio := float64(width) / float64(height)
	scale := 1.0 / float64(samplesPerPixel)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	return Image{
		width,
		height,
		samplesPerPixel,
		maxDepth,
		aspectRatio,
		scale,
		img,
	}
}

func (i Image) WriteColor(x, y int, c Color) {
	// Divide the c by the number of samples and gamma-correct for gamma 2.0
	r := math.Sqrt(c.X * i.Scale)
	g := math.Sqrt(c.Y * i.Scale)
	b := math.Sqrt(c.Z * i.Scale)

	// write the translated [0, 255] value of each color component
	rgba := color.RGBA{
		R: uint8(256 * Clamp(r, 0.0, 0.999)),
		G: uint8(256 * Clamp(g, 0.0, 0.999)),
		B: uint8(256 * Clamp(b, 0.0, 0.999)),
		A: 255,
	}

	i.image.Set(x, y, rgba)
}

func (i Image) Save(w io.Writer) error {
	return png.Encode(w, i.image)
}
