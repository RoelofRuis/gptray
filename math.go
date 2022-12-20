package main

import (
	"math"
	"math/rand"
)

func RandomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}
