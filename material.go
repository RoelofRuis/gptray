package main

import (
	"math"
)

type Material interface {
	Scatter(in Ray, rec HitRecord) (isScattered bool, attenuation Color, scatteredRay Ray)
	Emitted(p Vector, u, v float64) Color
}

type Lambertian struct {
	// Albedo is the surface color of the material
	Albedo Texture
}

func (l Lambertian) Scatter(in Ray, rec HitRecord) (bool, Color, Ray) {
	scatterDirection := rec.Normal.Add(RandomUnitVector())

	// Catch degenerate scatter direction
	if scatterDirection.NearZero() {
		scatterDirection = rec.Normal
	}

	scatteredRay := Ray{rec.P, scatterDirection, in.Wavelength}
	return true, l.Albedo.Value(rec.P, rec.U, rec.V), scatteredRay
}

func (l Lambertian) Emitted(p Vector, u, v float64) Color {
	return Color{}
}

type Metal struct {
	Albedo Color
	Fuzz   float64
}

func (m Metal) Scatter(in Ray, rec HitRecord) (bool, Color, Ray) {
	reflected := in.Direction.Unit().Reflect(rec.Normal)
	scatteredRay := Ray{rec.P, reflected.Add(RandomInUnitSphere().MulScalar(m.Fuzz)), in.Wavelength}
	return rec.Normal.Dot(scatteredRay.Direction) > 0, m.Albedo, scatteredRay
}

func (m Metal) Emitted(p Vector, u, v float64) Color {
	return Color{}
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

	scattered := Ray{rec.P, direction, in.Wavelength}
	return true, attenuation, scattered
}

func (d Dielectric) Emitted(p Vector, u, v float64) Color {
	return Color{}
}

func reflectance(cosine, refIdx float64) float64 {
	// Use Schlick's approximation for reflectance
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

type DiffuseLight struct {
	Emit Texture
}

func (d DiffuseLight) Scatter(in Ray, rec HitRecord) (isScattered bool, attenuation Color, scatteredRay Ray) {
	return false, Color{}, Ray{}
}

func (d DiffuseLight) Emitted(p Vector, u, v float64) Color {
	return d.Emit.Value(p, u, v)
}
