package main

import "math"

const MaxRecursionDepth = 2

func Shade(scene *Scene, ray *Ray, intersection *Intersection, depth int) *Color {
	// calculate the color contribution at the intersection point
	color := CalculateColor(scene, ray, intersection)

	// check if the maximum recursion depth has been reached
	if depth >= MaxRecursionDepth {
		return &color
	}

	// calculate the transparency and reflectivity of the material at the intersection point
	transparency := intersection.Object.Material().Transparency
	reflectivity := intersection.Object.Material().Reflectivity

	// calculate the reflection and refraction contributions
	if transparency > 0 || reflectivity > 0 {
		// calculate the reflection and refraction rays
		reflectionRay, refractionRay := CalculateReflectionAndRefractionRays(ray, intersection)

		reflectionIntersection, reflectionFound := FindIntersection(scene, reflectionRay)
		refractionIntersection, refractionFound := FindIntersection(scene, refractionRay)

		// calculate the reflection and refraction colors
		reflectionColor := Color{}
		if reflectionFound {
			reflectionColor = Shade(scene, reflectionRay, reflectionIntersection, depth+1).MultiplyScalar(reflectivity)
		}
		refractionColor := Color{}
		if refractionFound {
			refractionColor = Shade(scene, refractionRay, refractionIntersection, depth+1).MultiplyScalar(transparency)
		}

		// combine the reflection and refraction colors with the oclor contribution at the intersection point
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
	if intersection.Object.Material().Transparency > 0 {
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

	// get the diffuse color and intensity of the material at the intersection point
	diffuseColor := intersection.Object.Material().Color
	diffuseIntensity := intersection.Object.Material().DiffuseIntensity

	// get the normal vector at the intersection point
	normal := intersection.Object.Normal(&intersection.Point)

	// iterate over all the lights in the scene
	for _, light := range scene.Lights {
		// calculate the light direction vector
		lightDirection := light.Position.Subtract(intersection.Point).Normalize()

		// check if the light is blocked by another object
		lightRay := &Ray{intersection.Point, lightDirection}
		lightIntersection, found := FindIntersection(scene, lightRay)
		if found && lightIntersection.Distance < light.Position.Subtract(intersection.Point).Length() {
			continue
		}

		// calculate the diffuse lighting contribution
		diffuse = diffuse.Add(diffuseColor.MultiplyScalar(diffuseIntensity * math.Max(0, normal.Dot(lightDirection))))
	}

	return diffuse
}

// CalculateSpecularLighting calculates the specular lighting contribution at the intersection point.
func CalculateSpecularLighting(scene *Scene, ray *Ray, intersection *Intersection) Color {
	//initialize the specular lighting contribution to zero
	specular := Color{}

	// get the specular color, intensity, and shininess of the material at the intersection point
	specularColor := intersection.Object.Material().Color
	specularIntensity := intersection.Object.Material().SpecularIntensity
	shininess := intersection.Object.Material().Shininess

	// get the normal vector at the intersection point
	normal := intersection.Object.Normal(&intersection.Point)

	for _, light := range scene.Lights {
		// calculate the light direction vector
		lightDirection := light.Position.Subtract(intersection.Point).Normalize()

		// check if the light is blocked by another object
		lightRay := &Ray{intersection.Point, lightDirection}
		lightIntersection, found := FindIntersection(scene, lightRay)
		if found && lightIntersection.Distance < light.Position.Subtract(intersection.Point).Length() {
			continue
		}

		halfVector := ray.Direction.Add(lightDirection).Normalize()
		specular = specular.Add(specularColor.MultiplyScalar(specularIntensity * math.Pow(math.Max(0, normal.Dot(halfVector)), shininess)))
	}

	return specular
}

func Clamp(value float64, min float64, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	} else {
		return value
	}
}

