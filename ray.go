package main

// Ray represents a ray in 3D space
type Ray struct {
	Origin     Vector
	Direction  Vector
	Wavelength float64
}

// At returns the point at distance t along the ray
func (r Ray) At(t float64) Vector {
	return r.Origin.Add(r.Direction.MulScalar(t))
}
