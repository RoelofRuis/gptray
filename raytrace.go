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
			// calculate the x and y coordinates of the center of the pixel in screen space
			screenX := (float64(x) - halfWidth) / halfWidth
			screenY := -(float64(y) -halfHeight) / halfHeight

			// create a ray that passes through the center of the pixel
			ray := CreateRay(scene.Camera, screenX, screenY, aspectRatio)

			color := SampleColor(scene, ray, 5, 0)

			img.Set(x, y, color)
		}
	}

	return img
}

func SampleColor(scene *Scene, ray *Ray, sampleSize, depth int) Color {
	// initialize the color to black
	color := Color{}

	// create a grid of samples
	for sy := 0; sy < sampleSize; sy++ {
		for sx := 0; sx < sampleSize; sx++ {
			//adjust the position of the ray slightly for each sample
			sampleRay := &Ray{
				Origin: ray.Origin,
				Direction: ray.Direction.Add(ray.Direction.MultiplyScalar((float64(sx * 1) - float64(sampleSize - 1)/2) /
					float64(sampleSize)).Add(ray.Direction.MultiplyScalar((float64(sy * 1) - float64(sampleSize - 1)/2) / float64(sampleSize)))),
			}

			intersection, found := FindIntersection(scene, sampleRay)

			if !found {
				// if there is no intersection, set the color to the background color
				color = color.Add(*scene.AmbientColor)
				continue
			}

			// calculate the color based on the intersection point, surface normal, material properties, and lighting conditions
			shade := Shade(scene, sampleRay, intersection, depth)
			color = color.Add(*shade)
		}
	}

	color = color.MultiplyScalar(1.0 / float64(sampleSize*sampleSize))

	return color
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
