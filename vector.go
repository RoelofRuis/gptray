package main

import (
	"fmt"
	"math"
)

type Vector struct {
	X, Y, Z float64
}

// Neg returns a new Vector that is the negation of v
func (v Vector) Neg() Vector {
	return Vector{-v.X, -v.Y, -v.Z}
}

// Add returns a new Vector that is the sum of v and w.
func (v Vector) Add(w Vector) Vector {
	return Vector{v.X + w.X, v.Y + w.Y, v.Z + w.Z}
}

// Sub returns a new Vector that is the difference between v and w.
func (v Vector) Sub(w Vector) Vector {
	return v.Add(w.Neg())
}

func (v Vector) Mul(w Vector) Vector {
	return Vector{v.X * w.X, v.Y * w.Y, v.Z * w.Z}
}

// MulScalar returns a new Vector that is v scaled by s.
func (v Vector) MulScalar(s float64) Vector {
	return Vector{v.X * s, v.Y * s, v.Z * s}
}

// DivScalar returns a new Vector that is v scaled by 1/s.
func (v Vector) DivScalar(s float64) Vector {
	return v.MulScalar(1 / s)
}

// Length returns the length of v.
func (v Vector) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// LengthSquared returns the square of the length of v.
func (v Vector) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Unit returns a unit vector in the same direction as v.
func (v Vector) Unit() Vector {
	length := v.Length()
	if length == 0 {
		return Vector{}
	}
	return v.DivScalar(length)
}

// Dot returns the dot product of v and w.
func (v Vector) Dot(w Vector) float64 {
	return v.X*w.X + v.Y*w.Y + v.Z*w.Z
}

// Cross returns the cross product of v and w.
func (v Vector) Cross(w Vector) Vector {
	return Vector{v.Y*w.Z - v.Z*w.Y, v.Z*w.X - v.X*w.Z, v.X*w.Y - v.Y*w.X}
}

// String returns a string representation of u in the form "(x, y, z)".
func (v Vector) String() string {
	return fmt.Sprintf("(%g, %g, %g)", v.X, v.Y, v.Z)
}

// NearZero returns true if the vector is close to zero in all dimensions.
func (v Vector) NearZero() bool {
	e := 1e-8
	return math.Abs(v.X) < e && math.Abs(v.Y) < e && math.Abs(v.Z) < e
}

func (v Vector) Reflect(n Vector) Vector {
	return v.Sub(n.MulScalar(2 * v.Dot(n)))
}

func (v Vector) Refract(n Vector, etaIOverEtaT float64) Vector {
	cosTheta := math.Min(v.Neg().Dot(n), 1.0)
	rOutPerp := v.Add(n.MulScalar(cosTheta)).MulScalar(etaIOverEtaT)
	rOutParallel := n.MulScalar(-math.Sqrt(math.Abs(1.0 - rOutPerp.LengthSquared())))
	return rOutPerp.Add(rOutParallel)
}

func Random(min, max float64) Vector {
	return Vector{RandomFloat64(min, max), RandomFloat64(min, max), RandomFloat64(min, max)}
}

func RandomInUnitSphere() Vector {
	for {
		p := Random(-1.0, 1.0)
		if p.LengthSquared() >= 1 {
			continue
		}
		return p
	}
}

func RandomUnitVector() Vector {
	return RandomInUnitSphere().Unit()
}

func RandomInUnitDisk() Vector {
	for {
		p := Vector{RandomFloat64(-1, 1), RandomFloat64(-1, 1), 0}
		if p.LengthSquared() >= 1 {
			continue
		}
		return p
	}
}
