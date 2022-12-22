package main

import "math"

type World struct {
	Hittables  []Hittable
	Background Color
}

func (w World) Hit(r Ray, tMin, tMax float64) (HitRecord, bool) {
	var rec HitRecord
	hasHit := false
	closest := tMax

	for _, hittable := range w.Hittables {
		if hr, ok := hittable.Hit(r, tMin, closest); ok {
			hasHit = true
			closest = hr.T
			rec = hr
		}
	}

	return rec, hasHit
}

// HitRecord represents the data of a ray-object intersection
type HitRecord struct {
	P         Vector   // position of the intersection point of the ray and the object
	U         float64  // U surface coordinate of the ray-object hitpoint
	V         float64  // V surface coordinate of the ray-object hitpoint
	Normal    Vector   // surface normal at the intersection point
	T         float64  // the distance along the ray from the origin to the intersection
	FrontFace bool     // whether the ray hit the front or the back of the object
	Material  Material // the material that was hit
}

func (h *HitRecord) SetFaceNormal(r Ray, outwardNormal Vector) {
	h.FrontFace = r.Direction.Dot(outwardNormal) < 0
	if h.FrontFace {
		h.Normal = outwardNormal
	} else {
		h.Normal = outwardNormal.Neg()
	}
}

// Hittable is an interface for objects that can be intersected by a ray.
type Hittable interface {
	// Hit returns true if the given ray intersects the object, and it returns
	// the hit record.
	Hit(r Ray, tMin, tMax float64) (HitRecord, bool)
}

type Sphere2 struct {
	Center   Vector
	Radius   float64
	Material Material
}

func (s Sphere2) Hit(r Ray, tMin, tMax float64) (HitRecord, bool) {
	oc := r.Origin.Sub(s.Center)
	a := r.Direction.LengthSquared()
	halfB := oc.Dot(r.Direction)
	c := oc.LengthSquared() - s.Radius*s.Radius

	discriminant := halfB*halfB - a*c
	if discriminant < 0 {
		return HitRecord{}, false
	}

	sqrtd := math.Sqrt(discriminant)
	root := (-halfB - sqrtd) / a
	if root < tMin || tMax < root {
		root = (-halfB + sqrtd) / a
		if root < tMin || tMax < root {
			return HitRecord{}, false
		}
	}

	t := root
	p := r.At(t)
	outwardNormal := p.Sub(s.Center).DivScalar(s.Radius)

	theta := math.Acos(-p.Y)
	phi := math.Atan2(-p.Z, p.X) + math.Pi

	u := phi / (2 * math.Pi)
	v := theta / math.Pi

	hitRecord := HitRecord{P: p, U: u, V: v, T: t, Material: s.Material}
	hitRecord.SetFaceNormal(r, outwardNormal)

	return hitRecord, true
}
