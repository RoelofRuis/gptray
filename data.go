package main

import "math"

type Scene struct {
	Objects          []Object // slice of objects in the scene
	Lights           []Light  // slice of light sources in the scene
	Camera           Camera2  // camera position and orientation
	AmbientColor     *Color2  // ambient light color of the scene
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
	Color     *Color2 // color of the light source
}

// Deprecated
type Camera2 struct {
	Position Vector  // position of the camera
	LookAt   Vector  // point that the camera is looking at
	Up       Vector  // up direction of the camera
	Fov      float64 // field of view of the camera (in degrees)
}

type Material struct {
	Color            *Color2 //diffuse color of the material
	Diffuse          float64 // diffuse lighting coefficient (0-1)
	Specular         float64 // specular lighting coefficient (0-1)
	SpecularExponent float64 // shininess coefficient (0-infinity)
	Reflective       float64 // reflective coefficient (0-1)
	Refractive       float64 // refractive coefficient (0-1)
	RefractionIndex  float64 // index of refraction (1-infinity)
}

// Deprecated
type Color2 struct {
	R, G, B float64
}

func (c Color2) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R * 65535)
	g = uint32(c.G * 65535)
	b = uint32(c.B * 65535)
	a = 65535
	return
}

func (c Color2) Add(c2 Color2) Color2 {
	return Color2{
		R: c.R + c2.R,
		G: c.G + c2.G,
		B: c.B + c2.B,
	}
}

func (c Color2) Clamp() Color2 {
	return Color2{
		math.Max(0, math.Min(c.R, 1)),
		math.Max(0, math.Min(c.G, 1)),
		math.Max(0, math.Min(c.B, 1)),
	}
}

func (c Color2) MultiplyScalar(s float64) Color2 {
	return Color2{
		R: c.R * s,
		G: c.G * s,
		B: c.B * s,
	}
}

func (c Color2) Multiply(c2 Color2) Color2 {
	return Color2{c.R * c2.R, c.G * c2.G, c.B * c2.B}
}

type Intersection struct {
	Object   Object  // the object that was intersected
	Point    Vector  // the intersection point
	Distance float64 // the distance from the ray origin to the intersection point
	Inside   bool    // whether the ray is inside the object at the intersection point
}
