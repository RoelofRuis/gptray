package main

import (
	"math"
)

type Material interface {
	Scatter(in Ray, rec HitRecord) (isScattered bool, attenuation Color, scatteredRay Ray)
}

type Lambertian struct {
	// Albedo is the surface color of the material
	Albedo Color
}

func (l Lambertian) Scatter(in Ray, rec HitRecord) (bool, Color, Ray) {
	scatterDirection := rec.Normal.Add(RandomUnitVector())

	// Catch degenerate scatter direction
	if scatterDirection.NearZero() {
		scatterDirection = rec.Normal
	}

	scatteredRay := Ray{rec.Position, scatterDirection}
	return true, l.Albedo, scatteredRay
}

type Metal struct {
	Albedo Color
	Fuzz   float64
}

func (m Metal) Scatter(in Ray, rec HitRecord) (bool, Color, Ray) {
	reflected := in.Direction.Unit().Reflect(rec.Normal)
	scatteredRay := Ray{rec.Position, reflected.Add(RandomInUnitSphere().MulScalar(m.Fuzz))}
	return rec.Normal.Dot(scatteredRay.Direction) > 0, m.Albedo, scatteredRay
}

type Dielectric struct {
	// Ir is the index of refraction of the material
	Ir float64
}

func (d Dielectric) Scatter(in Ray, rec HitRecord) (bool, Color, Ray) {
	attenuation := Color{1, 1, 1}
	refractionRatio := d.Ir
	if rec.FrontFace {
		refractionRatio = 1.0 / d.Ir
	}

	unitDirection := in.Direction.Unit()
	cosTheta := math.Min(unitDirection.Neg().Dot(rec.Normal), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)

	cannotRefract := refractionRatio*sinTheta > 1.0
	var direction Vector
	if cannotRefract || reflectance(cosTheta, refractionRatio) > RandomFloat64(0, 1) {
		direction = unitDirection.Reflect(rec.Normal)
	} else {
		direction = unitDirection.Refract(rec.Normal, refractionRatio)
	}

	scattered := Ray{rec.Position, direction}
	return true, attenuation, scattered
}

func reflectance(cosine, refIdx float64) float64 {
	// Use Schlick's approximation for reflectance
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}
