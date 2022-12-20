package main

import (
	"image"
	"math"
	"math/rand"
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
			// calculate the x and y coordinates of the center of the pixel in screen space
			screenX := (float64(x) - halfWidth) / halfWidth
			screenY := -(float64(y) - halfHeight) / halfHeight

			// create a ray that passes through the center of the pixel
			ray := CreateRay(scene.Camera, screenX, screenY, aspectRatio)

			color := SampleColor(scene, ray, SampleSize, 0)

			img.Set(x, y, color)
		}
	}

	return img
}

func SampleColor(scene *Scene, ray *Ray, sampleSize, depth int) Color2 {
	// initialize the color to black
	color := Color2{}

	offset := 0.01

	for n := 0; n < sampleSize; n++ {
		x := (rand.Float64() * 2 * offset) - offset
		y := (rand.Float64() * 2 * offset) - offset
		z := (rand.Float64() * 2 * offset) - offset

		sampleRay := &Ray{
			Origin:    ray.Origin,
			Direction: ray.Direction.Add(Vector{x, y, z}),
		}

		intersection, found := FindIntersection(scene, sampleRay)

		if !found {
			// if there is no intersection, set the color to the background color
			color = color.Add(*scene.AmbientColor)
			continue
		}

		// calculate the color based on the intersection point, surface normal, material properties, and lighting conditions
		shade := Shade(scene, sampleRay, intersection, depth)
		color = color.Add(shade.Clamp())
	}

	color = color.MultiplyScalar(1.0 / float64(sampleSize))

	return color
}

func CreateRay(camera Camera2, screenX, screenY, aspectRatio float64) *Ray {
	// Calculate the FOV angle in radians
	fovRadians := camera.Fov * math.Pi / 180.0

	// Calculate hte distance of the near plane from the camera
	nearPlaneDistance := 1.0 / math.Tan(fovRadians/2.0)

	// Calculate the direction of the ray
	dir := camera.LookAt.Sub(camera.Position).Unit()

	// Calculate the right and up vectors
	right := camera.Up.Cross(dir).Unit()
	up := dir.Cross(right)

	// Scale the right and up vectors by the screen coordinates
	right = right.MulScalar(screenX * aspectRatio)
	up = up.MulScalar(screenY)

	// Add the right and up vectors to the direction vector
	dir = dir.Add(right).Add(up).Unit()

	// Scale the direction vector by the near plane distance
	dir = dir.MulScalar(nearPlaneDistance)

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
			intersectionPoint := ray.Origin.Add(ray.Direction.MulScalar(distance))

			inside := false
			if intersectionPoint.Sub(ray.Origin).Dot(ray.Direction) < 0 {
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
