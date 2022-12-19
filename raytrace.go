package main

import (
	"image"
	"math"
)

func Raytrace(scene *Scene, width, height int) *image.RGBA {
	// create an image with the given width and height
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// calculate the aspect ratio of the image
	aspectRatio := float64(width) / float64(height)

	// calculate the half width and half height of the image
	halfWidth := float64(width) / 2.0
	halfHeight := float64(height) / 2.0

	// iterate over all the pixels in the image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// initialize the pixel color to black
			color := Color{0, 0, 0}

			// create a grid of samples for the pixel
			sampleSize := 5 // number of samples per dimension (e.g. 2x2, 3x3, etc.)
			for sy := 0; sy < sampleSize; sy++ {
				for sx := 0; sx < sampleSize; sx++ {
					// calculate the x and y coordinates of the sample in screen sapce
					screenX := (float64(x) - halfWidth + (float64(sx)+0.5)/float64(sampleSize)) / halfWidth
					screenY := -(float64(y) - halfHeight + (float64(sy)+0.5)/float64(sampleSize)) / halfHeight

					// create a ray that passes through the sample
					ray := CreateRay(scene.Camera, screenX, screenY, aspectRatio)

					intersection, found := FindIntersection(scene, ray)

					if !found {
						// if there is no intersection, set the pixel color to the background color
						color = color.Add(*scene.BackgroundColor)
						continue
					}

					// calculate the color of the pixel based on the intersection point, surface normal, material properties, and lighting conditions
					color = color.Add(*Shade(scene, ray, intersection, 0))
				}
			}

			// average the colors of the samples
			color = color.MultiplyScalar(1.0 / float64(sampleSize*sampleSize))

			img.Set(x, y, color)
		}
	}

	return img
}

func CreateRay(camera Camera, screenX, screenY, aspectRatio float64) *Ray {
	// Calculate the FOV angle in radians
	fovRadians := camera.Fov * math.Pi / 180.0

	// Calculate hte distance of the near plane from the camera
	nearPlaneDistance := 1.0 / math.Tan(fovRadians/2.0)

	// Calculate the direction of the ray
	dir := camera.LookAt.Subtract(camera.Position).Normalize()

	// Calculate the right and up vectors
	right := camera.Up.Cross(dir).Normalize()
	up := dir.Cross(right)

	// Scale the right and up vectors by the screen coordinates
	right = right.MultiplyScalar(screenX * aspectRatio)
	up = up.MultiplyScalar(screenY)

	// Add the right and up vectors to the direction vector
	dir = dir.Add(right).Add(up).Normalize()

	// Scale the direction vector by the near plane distance
	dir = dir.MultiplyScalar(nearPlaneDistance)

	// return the ray with the given origin and direction
	return &Ray{
		Origin:    camera.Position,
		Direction: dir,
	}
}

func FindIntersection(scene *Scene, ray *Ray) (*Intersection, bool) {
	// initialize the nearest intersection to a large value
	nearest := math.MaxFloat64
	var nearestIntersection *Intersection

	// iterate over all the objects in the scene
	for _, object := range scene.Objects {
		// calculate the intersection of the ray with the object
		distance, found := object.Intersect(ray)
		if !found {
			continue
		}

		// check if the intersection is the nearest so far
		if distance < nearest {
			// calculate the intersection point
			intersectionPoint := ray.Origin.Add(ray.Direction.MultiplyScalar(distance))

			inside := false
			if intersectionPoint.Subtract(ray.Origin).Dot(ray.Direction) < 0 {
				inside = true
			}

			// update the nearest intersection
			nearest = distance
			nearestIntersection = &Intersection{
				Point:    intersectionPoint,
				Distance: distance,
				Object:   object,
				Inside:   inside,
			}
		}
	}

	if nearestIntersection == nil {
		// if no intersection was found, return false
		return nil, false
	}

	//return the nearest intersection
	return nearestIntersection, true
}

const MAX_RECURSION_DEPTH = 1

func Shade(scene *Scene, ray *Ray, intersection *Intersection, depth int) *Color {
	// Stop recursion when maximum depth is reached
	if depth > MAX_RECURSION_DEPTH {
		return &Color{}
	}

	point := intersection.Point
	normal := intersection.Object.Normal(&point)
	material := intersection.Object.Material()

	// Calculate ambient light contribution
	color := material.Color.MultiplyScalar(scene.AmbientLight)

	// Calculate diffuse and specular light contributions
	for _, light := range scene.Lights {
		// Calculate direction from point to light source
		lightDirection := light.Position.Subtract(point).Normalize()

		// Calculate diffuse light contribution
		diffuse := lightDirection.Dot(normal)
		if diffuse > 0 {
			diffuseColor := material.Color.MultiplyScalar(diffuse)
			diffuseColor = diffuseColor.MultiplyScalar(light.Intensity)
			color = color.Add(diffuseColor)
		}

		// Calculate specular light contribution
		if material.Specular > 0 {
			viewDirection := scene.Camera.Position.Subtract(point).Normalize()
			halfVector := lightDirection.Add(viewDirection).Normalize()
			specular := math.Pow(halfVector.Dot(normal), material.Shininess)
			if specular > 0 {
				specularColor := light.Color.MultiplyScalar(specular)
				specularColor = specularColor.MultiplyScalar(material.Specular * light.Intensity)
				color = color.Add(specularColor)
			}
		}
	}

	// Calculate reflective contribution
	if material.Reflective > 0 {
		reflectionDirection := ray.Direction.Subtract(normal.MultiplyScalar(2 * ray.Direction.Dot(normal))).Normalize()
		reflectionRay := &Ray{point, reflectionDirection}
		reflectionIntersection, found := FindIntersection(scene, reflectionRay)
		if found {
			reflectionColor := Shade(scene, reflectionRay, reflectionIntersection, depth+1).MultiplyScalar(material.Reflective)
			color = color.Add(reflectionColor)
		}
	}

	// Calculate transparent and refractive contributions
	if material.Transparent > 0 {
		refractionIndex := material.Refraction
		if intersection.Inside {
			refractionIndex = 1.0 / refractionIndex
		}
		cosI := -normal.Dot(ray.Direction)
		sinT2 := refractionIndex * refractionIndex * (1.0 - cosI*cosI)
		if sinT2 > 1.0 {
			// Total internal reflection
			reflectionDirection := ray.Direction.Subtract(normal.MultiplyScalar(2 * ray.Direction.Dot(normal))).Normalize()
			reflectionRay := &Ray{point, reflectionDirection}
			reflectionIntersection, found := FindIntersection(scene, reflectionRay)
			if found {
				reflectionColor := Shade(scene, reflectionRay, reflectionIntersection, depth+1).MultiplyScalar(material.Transparent)
				color = color.Add(reflectionColor)
			}
		} else {
			// Refraction
			cosT := math.Sqrt(1.0 - sinT2)
			refractionDirection := ray.Direction.MultiplyScalar(refractionIndex).Add(normal.MultiplyScalar(refractionIndex*cosI - cosT)).Normalize()
			refractionRay := &Ray{point, refractionDirection}
			refractionIntersection, found := FindIntersection(scene, refractionRay)
			if found {
				refractionColor := Shade(scene, refractionRay, refractionIntersection, depth+1).MultiplyScalar(material.Transparent)
				color = color.Add(refractionColor)
			}
		}
	}

	return &color
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
