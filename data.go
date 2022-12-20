package main

import "math"

type Scene struct {
	Objects          []Object // slice of objects in the scene
	Lights           []Light  // slice of light sources in the scene
	Camera           Camera   // camera position and orientation
	AmbientColor     *Color   // ambient light color of the scene
	AmbientIntensity float64  // ambient light intensity in the scene
}

type Object interface {
	// Intersect calculates the intersection point of a ray with the object
	Intersect(ray *Ray) (float64, bool)
	// Normal calculates the surface normal at a given point on the object
	Normal(point *Vector) Vector
	// Material properties of the object
	Material() Material
}

type Light struct {
	Position  Vector  // position of the light source
	Intensity float64 // intensity of the light source
	Color     *Color  // color of the light source
}

type Camera struct {
	Position Vector  // position of the camera
	LookAt   Vector  // point that the camera is looking at
	Up       Vector  // up direction of the camera
	Fov      float64 // field of view of the camera (in degrees)
}

type Material struct {
	Color            *Color  //diffuse color of the material
	Diffuse          float64 // diffuse lighting coefficient (0-1)
	Specular         float64 // specular lighting coefficient (0-1)
	SpecularExponent float64 // shininess coefficient (0-infinity)
	Reflective       float64 // reflective coefficient (0-1)
	Refractive       float64 // refractive coefficient (0-1)
	RefractionIndex  float64 // index of refraction (1-infinity)
}

type Color struct {
	R, G, B float64
}

func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R * 65535)
	g = uint32(c.G * 65535)
	b = uint32(c.B * 65535)
	a = 65535
	return
}

func (c Color) Add(c2 Color) Color {
	return Color{
		R: c.R + c2.R,
		G: c.G + c2.G,
		B: c.B + c2.B,
	}
}

func (c Color) Clamp() Color {
	return Color{
		math.Max(0, math.Min(c.R, 1)),
		math.Max(0, math.Min(c.G, 1)),
		math.Max(0, math.Min(c.B, 1)),
	}
}

func (c Color) MultiplyScalar(s float64) Color {
	return Color{
		R: c.R * s,
		G: c.G * s,
		B: c.B * s,
	}
}

func (c Color) Multiply(c2 Color) Color {
	return Color{c.R * c2.R, c.G * c2.G, c.B * c2.B}
}

type Vector struct {
	X, Y, Z float64
}

func (v Vector) Add(w Vector) Vector {
	return Vector{
		X: v.X + w.X,
		Y: v.Y + w.Y,
		Z: v.Z + w.Z,
	}
}

func (v Vector) Subtract(w Vector) Vector {
	return Vector{
		X: v.X - w.X,
		Y: v.Y - w.Y,
		Z: v.Z - w.Z,
	}
}

func (v Vector) MultiplyScalar(s float64) Vector {
	return Vector{
		X: v.X * s,
		Y: v.Y * s,
		Z: v.Z * s,
	}
}

func (v Vector) Normalize() Vector {
	length := v.Length()
	if length == 0 {
		return Vector{}
	}
	return Vector{
		X: v.X / length,
		Y: v.Y / length,
		Z: v.Z / length,
	}
}

func (v Vector) DirectionTo(w Vector) Vector {
	displacement := w.Subtract(v)
	return displacement.Normalize()
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector) Dot(w Vector) float64 {
	return v.X*w.X + v.Y*w.Y + v.Z*w.Z
}

func (v Vector) Cross(w Vector) Vector {
	return Vector{
		X: v.Y*w.Z - v.Z*w.Y,
		Y: v.Z*w.X - v.X*w.Z,
		Z: v.X*w.Y - v.Y*w.X,
	}
}

func (v Vector) Reflect(n Vector) Vector {
	return v.Subtract(n.MultiplyScalar(2 * v.Dot(n)))
}

func (v Vector) Negate() Vector {
	return Vector{-v.X, -v.Y, -v.Z}
}

type Ray struct {
	Origin    Vector // origin of the ray
	Direction Vector // direction of the ray
}

type Intersection struct {
	Object   Object  // the object that was intersected
	Point    Vector  // the intersection point
	Distance float64 // the distance from the ray origin to the intersection point
	Inside   bool    // whether the ray is inside the object at the intersection point
}
