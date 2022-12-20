package main

type Camera struct {
	Origin          Vector
	lowerLeftCorner Vector
	horizontal      Vector
	vertical        Vector
}

func NewCamera(aspectRatio float64) *Camera {
	viewportHeight := 2.0
	viewportWidth := aspectRatio * viewportHeight
	focalLength := 1.0

	origin := Vector{0, 0, 0}
	horizontal := Vector{viewportWidth, 0, 0}
	vertical := Vector{0, viewportHeight, 0}
	lowerLeftCorner := origin.
		Sub(horizontal.MulScalar(0.5)).
		Sub(vertical.MulScalar(0.5)).
		Sub(Vector{0, 0, focalLength})

	return &Camera{origin, lowerLeftCorner, horizontal, vertical}
}

func (c Camera) GetRay(x, y float64) Ray {
	return Ray{
		c.Origin,
		c.lowerLeftCorner.Add(c.horizontal.MulScalar(x)).Add(c.vertical.MulScalar(y)).Sub(c.Origin),
	}
}
