package main

type Plane struct {
	Position           Vector   // position of the plane
	NormalVector       Vector   // normal of the plane
	MaterialProperties Material // material properties of the plane
}

func (p Plane) Intersect(ray *Ray) (float64, bool) {
	// calculate the dot product of the normal and the ray direction
	dot := p.NormalVector.Dot(ray.Direction)

	// check if the ray is parallel to the plane
	if dot == 0 {
		// The ray is parallel to the plane, so it does not intersect
		return 0, false
	}

	// calculate the distance from the ray origin to the plane
	distance := p.NormalVector.Dot(p.Position.Subtract(ray.Origin)) / dot

	// Check if the distance is positive (i.e. the intersection point is in front of the ray)
	if distance <= 0 {
		// the intersection point is behind the ray, so it does not intersect
		return 0, false
	}

	// The ray intersects the plane at the given distance
	return distance, true
}

func (p Plane) Normal(point *Vector) Vector {
	// the normal of the plane is constant and does not depend on the point
	return p.NormalVector
}

func (p Plane) Material() Material {
	// Return the material properties of the plane
	return p.MaterialProperties
}
