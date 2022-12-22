package main

import "math"

type Texture interface {
	Value(p Vector, u, v float64) Color
}

type SolidColor struct {
	ColorValue Color
}

func NewSolidColor(r, g, b float64) SolidColor {
	return SolidColor{Color{r, g, b}}
}

func (s SolidColor) Value(p Vector, u, v float64) Color {
	return s.ColorValue
}

type CheckerTexture struct {
	Odd  Texture
	Even Texture
}

func (s CheckerTexture) Value(p Vector, u, v float64) Color {
	sines := math.Sin(10*p.X) * math.Sin(10*p.Y) * math.Sin(10*p.Z)
	if sines < 0 {
		return s.Odd.Value(p, u, v)
	} else {
		return s.Even.Value(p, u, v)
	}
}
