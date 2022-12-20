package main

import (
	"math"
)

const MaxRecursionDepth = 1
const SampleSize = 20

func Shade(scene *Scene, ray *Ray, intersection *Intersection, depth int) *Color {
	// calculate the color contribution at the intersection point
	color := CalculateColor(scene, ray, intersection)

	// check if the maximum recursion depth has been reached
	if depth >= MaxRecursionDepth {
		return &color
	}

	// calculate the transparency and reflectivity of the material at the intersection point
	transparency := intersection.Object.Material().Refractive
	reflectivity := intersection.Object.Material().Reflective

	// calculate the reflection and refraction contributions
	if transparency > 0 || reflectivity > 0 {
		// calculate the reflection and refraction rays
		reflectionRay, refractionRay := CalculateReflectionAndRefractionRays(ray, intersection)

		_, reflectionFound := FindIntersection(scene, reflectionRay)
		_, refractionFound := FindIntersection(scene, refractionRay)

		// calculate the reflection and refraction colors
		reflectionColor := Color{}
		if reflectionFound {
			reflectionColor = SampleColor(scene, reflectionRay, SampleSize, depth+1).MultiplyScalar(reflectivity)
		}
		refractionColor := Color{}
		if refractionFound {
			refractionColor = SampleColor(scene, refractionRay, SampleSize, depth+1).MultiplyScalar(transparency)
		}

		// combine the reflection and refraction colors with the color contribution at the intersection point
		color = color.MultiplyScalar(1 - transparency - reflectivity)
		color = color.Add(reflectionColor).Add(refractionColor)
	}

	return &color
}

// CalculateColor calculates the color contribution at the intersection point.
func CalculateColor(scene *Scene, ray *Ray, intersection *Intersection) Color {
	// calculate the ambient, diffuse, and specular lighting contributions
	ambient := CalculateAmbientLighting(scene, intersection)
	diffuse := CalculateDiffuseLighting(scene, ray, intersection)
	specular := CalculateSpecularLighting(scene, ray, intersection)

	color := ambient.Add(diffuse).Add(specular)

	return color
}

func CalculateReflectionAndRefractionRays(ray *Ray, intersection *Intersection) (*Ray, *Ray) {
	// calculate the normal vector at the intersection point
	normal := intersection.Object.Normal(&intersection.Point)

	reflectionDirection := ray.Direction.Subtract(normal.MultiplyScalar(2 * ray.Direction.Dot(normal))).Normalize()
	reflectionRay := &Ray{intersection.Point, reflectionDirection}

	refractionRay := &Ray{}
	if intersection.Object.Material().Refractive > 0 {
		refractionIndex := intersection.Object.Material().RefractionIndex
		if intersection.Inside {
			refractionIndex = 1.0 / refractionIndex
		}
		cosI := -normal.Dot(ray.Direction)
		sinT2 := refractionIndex * refractionIndex * (1.0 - cosI*cosI)
		if sinT2 <= 1.0 {
			cosT := math.Sqrt(1.0 - sinT2)
			refractionDirection := ray.Direction.MultiplyScalar(refractionIndex).Add(normal.MultiplyScalar(refractionIndex*cosI - cosT)).Normalize()
			refractionRay = &Ray{intersection.Point, refractionDirection}
		}
	}

	return reflectionRay, refractionRay
}

// CalculateAmbientLighting calculates the ambient lighting contribution at the intersection point.
func CalculateAmbientLighting(scene *Scene, intersection *Intersection) Color {
	// get the ambient color and intensity of the scene
	ambientColor := scene.AmbientColor
	ambientIntensity := scene.AmbientIntensity

	// calculate the ambient lighting contribution
	ambient := ambientColor.MultiplyScalar(ambientIntensity)

	return ambient
}

// CalculateDiffuseLighting calculates the diffuse lighting contribution at the intersection point.
func CalculateDiffuseLighting(scene *Scene, ray *Ray, intersection *Intersection) Color {
	// initialize the diffuse lighting contribution to zero
	diffuse := Color{}

	for _, light := range scene.Lights {
		// calculate the direction of the light
		lightDirection := intersection.Point.DirectionTo(light.Position)

		// Cast a shadow ray from the intersection point towards the light
		shadowRay := &Ray{intersection.Point, lightDirection}

		_, intersects := FindIntersection(scene, shadowRay)

		if !intersects {
			// Calculate the dot product of the light direction and the surface normal
			dot := intersection.Object.Normal(&intersection.Point).Dot(lightDirection)

			// calculate the intensity of the diffuse lighting
			intensity := light.Intensity * dot

			// Calculate the diffuse lighting color
			color := light.Color.Multiply(*intersection.Object.Material().Color)

			diffuse = diffuse.Add(color.MultiplyScalar(intersection.Object.Material().Diffuse).MultiplyScalar(intensity))
		}
	}

	return diffuse
}

// CalculateSpecularLighting calculates the specular lighting contribution at the intersection point.
func CalculateSpecularLighting(scene *Scene, ray *Ray, intersection *Intersection) Color {
	//initialize the specular lighting contribution to zero
	specular := Color{}

	// Get the material properties of the intersected object
	material := intersection.Object.Material()

	for _, light := range scene.Lights {
		// calculate the light direction vector
		lightDirection := light.Position.DirectionTo(intersection.Point)

		reflection := lightDirection.Reflect(intersection.Object.Normal(&intersection.Point))

		// Calculate the intensity of the specular lighting
		intensity := light.Intensity * math.Pow(reflection.Dot(ray.Direction.Negate()), material.SpecularExponent)

		// Calculate the specular lighting color
		color := light.Color.Multiply(*material.Color)

		specular = specular.Add(color.MultiplyScalar(material.Specular).MultiplyScalar(intensity))
	}

	return specular
}
