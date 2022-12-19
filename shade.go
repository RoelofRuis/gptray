package main

import (
	"math"
)

const MaxRecursionDepth = 1

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
			reflectionColor = SampleColor(scene, reflectionRay, 5, depth+1).MultiplyScalar(reflectivity)
		}
		refractionColor := Color{}
		if refractionFound {
			refractionColor = SampleColor(scene, refractionRay, 5, depth+1).MultiplyScalar(transparency)
		}

		// combine the reflection and refraction colors with the oclor contribution at the intersection point
		color = color.MultiplyScalar(1 - transparency - reflectivity)
		color = color.Add(reflectionColor).Add(refractionColor).Clamp()
	}

	return &color
}

// CalculateColor calculates the color contribution at the intersection point.
func CalculateColor(scene *Scene, ray *Ray, intersection *Intersection) Color {
	// calculate the ambient, diffuse, and specular lighting contributions
	ambient := CalculateAmbientLighting(scene, intersection)
	diffuse := CalculateDiffuseLighting(scene, ray, intersection)
	specular := CalculateSpecularLighting(scene, ray, intersection)

	color := ambient.Add(diffuse).Add(specular).Clamp()

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
		sinT2 := refractionIndex * refractionIndex * (1.0 - cosI *cosI)
		if sinT2 <= 1.0 {
			cosT := math.Sqrt(1.0 - sinT2)
			refractionDirection := ray.Direction.MultiplyScalar(refractionIndex).Add(normal.MultiplyScalar(refractionIndex*cosI-cosT)).Normalize()
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

	// Get the material properties of the intersected object
	material := intersection.Object.Material()

	// get the normal vector at the intersection point
	normal := intersection.Object.Normal(&intersection.Point)

	// iterate over all the lights in the scene
	for _, light := range scene.Lights {
		// calculate the light direction vector
		lightDirection := intersection.Point.DirectionTo(light.Position)

		// check if the light is blocked by another object
		lightRay := &Ray{intersection.Point, lightDirection}
		lightIntersection, found := FindIntersection(scene, lightRay)
		if found && lightIntersection.Distance < light.Position.Subtract(intersection.Point).Length() {
			continue
		}

		// Calculate the diffuse lighting contribution using the dot product between the normal and the light direction
		diffuseContribution := normal.Dot(lightDirection)

		// Clamp the diffuse contribution to 0 if it is negative
		if diffuseContribution < 0 {
			diffuseContribution = 0
		}

		diffuse = diffuse.Add(light.Color.MultiplyScalar(diffuseContribution * material.Diffuse).Multiply(*material.Color))
	}

	return diffuse
}

// CalculateSpecularLighting calculates the specular lighting contribution at the intersection point.
func CalculateSpecularLighting(scene *Scene, ray *Ray, intersection *Intersection) Color {
	//initialize the specular lighting contribution to zero
	specular := Color{}

	// Get the material properties of the intersected object
	material := intersection.Object.Material()

	// get the normal vector at the intersection point
	normal := intersection.Object.Normal(&intersection.Point)

	// Calculate the reflection vector using the formula R = 2(N.L)N - L, where N is the surface normal and L is the light direction
	reflection := normal.MultiplyScalar(2 * normal.Dot(ray.Direction)).Subtract(ray.Direction)

	for _, light := range scene.Lights {
		// calculate the light direction vector
		lightDirection := light.Position.DirectionTo(intersection.Point)

		// check if the light is blocked by another object
		lightRay := &Ray{intersection.Point, lightDirection}
		lightIntersection, found := FindIntersection(scene, lightRay)
		if found && lightIntersection.Distance < light.Position.Subtract(intersection.Point).Length() {
			continue
		}

		// Calculate the specular lighting contribution using the dot product between the reflection vector and the light direction
		specularContribution := reflection.Dot(lightDirection)

		// Clamp the specular contribution to 0 if it is negative or if the angle between the reflection vector and the light direction is greater than 90 degrees
		if specularContribution < 0 || math.Abs(specularContribution) < math.Pi/2 {
			specularContribution = 0
		}

		// Raise the specular contribution to the power of the material's specular exponent
		specularContribution = math.Pow(specularContribution, material.SpecularExponent)


		// Calculate the specular lighting color by multiplying the specular contribution by the light color and the material's specular color
		specular = specular.Add(light.Color.MultiplyScalar(specularContribution * material.Specular).Multiply(*material.Color))
	}

	return specular
}
