package main

import "math"

type Sphere struct {
	Center             Vector   // center of the sphere
	Radius             float64  // radius of the sphere
	MaterialProperties Material // material properties of the sphere
}

// Intersect calculates the intersection point of a ray with the sphere.
// It returns the distance from the ray origin to the intersection point,
// and a boolean value indicating whether an intersection occurred.
func (s *Sphere) Intersect(ray *Ray) (float64, bool) {
	// calculate the intersection using the quadratic formula
	oc := ray.Origin.Subtract(s.Center)
	a := ray.Direction.Dot(ray.Direction)
	b := 2 * oc.Dot(ray.Direction)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		// no intersection
		return 0, false
	}

	sqrtDiscriminant := math.Sqrt(discriminant)
	t1 := (-b + sqrtDiscriminant) / (2 * a)
	t2 := (-b - sqrtDiscriminant) / (2 * a)

	if t1 > t2 {
		t1, t2 = t2, t1
	}

	if t1 < 0 {
		// intersection is behind the ray origin
		t1 = t2
	}

	if t1 < 0 {
		// intersection is in front of the ray origin
		return 0, false
	}

	// intersection is in front of the ray origin
	return t1, true
}

// Normal calculates the surface normal at a given point on the sphere.
func (s *Sphere) Normal(point *Vector) Vector {
	// the surface normal is simply the normalized vector from the center of the sphere to the point
	return point.Subtract(s.Center).Normalize()
}

// Material returns the material properties of the sphere.
func (s *Sphere) Material() Material {
	return s.MaterialProperties
}
