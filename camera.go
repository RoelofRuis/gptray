package main

import "math"

type Camera struct {
	Origin          Vector
	LowerLeftCorner Vector
	Horizontal      Vector
	Vertical        Vector
	U               Vector
	V               Vector
	W               Vector
	LensRadius      float64
}

// NewCamera requires the vertical field-of-view in degrees and the aspect ratio and returns a camera.
func NewCamera(
	lookFrom Vector,
	lookAt Vector,
	up Vector,
	vfov, // vertical field-of-view in degrees
	aspectRatio float64,
	aperture float64,
	focusDisc float64,
) *Camera {
	theta := DegreesToRadians(vfov)
	h := math.Tan(theta / 2)
	viewportHeight := 2.0 * h
	viewportWidth := aspectRatio * viewportHeight

	w := lookFrom.Sub(lookAt).Unit()
	u := up.Cross(w).Unit()
	v := w.Cross(u)

	origin := lookFrom
	horizontal := u.MulScalar(focusDisc * viewportWidth)
	vertical := v.MulScalar(focusDisc * viewportHeight)
	lowerLeftCorner := origin.
		Sub(horizontal.MulScalar(0.5)).
		Sub(vertical.MulScalar(0.5)).
		Sub(w.MulScalar(focusDisc))

	lensRadius := aperture / 2

	return &Camera{
		origin,
		lowerLeftCorner,
		horizontal,
		vertical,
		u,
		v,
		w,
		lensRadius,
	}
}

func (c Camera) GetRay(x, y float64) Ray {
	rd := RandomInUnitDisk().MulScalar(c.LensRadius)
	offset := c.U.MulScalar(rd.X).Add(c.V.MulScalar(rd.Y))

	return Ray{
		c.Origin.Add(offset),
		c.LowerLeftCorner.Add(c.Horizontal.MulScalar(x)).Add(c.Vertical.MulScalar(y)).Sub(c.Origin).Sub(offset),
	}
}
