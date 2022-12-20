package main

import "math"

type Camera struct {
	Origin          Vector
	LowerLeftCorner Vector
	Horizontal      Vector
	Vertical        Vector
}

// NewCamera requires the vertical field-of-view in degrees and the aspect ratio and returns a camera.
func NewCamera(
	lookFrom Vector,
	lookAt Vector,
	up Vector,
	vfov, // vertical field-of-view in degrees
	aspectRatio float64,
) *Camera {
	theta := DegreesToRadians(vfov)
	h := math.Tan(theta / 2)
	viewportHeight := 2.0 * h
	viewportWidth := aspectRatio * viewportHeight

	w := lookFrom.Sub(lookAt).Unit()
	u := up.Cross(w).Unit()
	v := w.Cross(u)

	origin := lookFrom
	horizontal := u.MulScalar(viewportWidth)
	vertical := v.MulScalar(viewportHeight)
	lowerLeftCorner := origin.
		Sub(horizontal.MulScalar(0.5)).
		Sub(vertical.MulScalar(0.5)).
		Sub(w)

	return &Camera{origin, lowerLeftCorner, horizontal, vertical}
}

func (c Camera) GetRay(x, y float64) Ray {
	return Ray{
		c.Origin,
		c.LowerLeftCorner.Add(c.Horizontal.MulScalar(x)).Add(c.Vertical.MulScalar(y)).Sub(c.Origin),
	}
}
