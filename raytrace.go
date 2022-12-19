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
			//calculate the x and y coordinates of the pixel in screen space
			screenX := (float64(x) - halfWidth) / halfWidth
			screenY := -(float64(y) - halfHeight) / halfHeight

			// create a ray that passes through the pixel
			ray := CreateRay(scene.Camera, screenX, screenY, aspectRatio)

			intersection, found := FindIntersection(scene, ray)
			if !found {
				// if there is no intersection, set the pixel color to the background color
				img.Set(x, y, scene.BackgroundColor)
				continue
			}

			// calculate the color of the pixel based on the intersection point, surface normal, material properties, and lighting conditions
			color := Shade(scene, ray, intersection)

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
	dir = dir.Add(camera.Up.MultiplyScalar(screenY)).Normalize()
	dir = dir.Add(camera.Up.Cross(dir).Normalize().MultiplyScalar(screenX / aspectRatio)).Normalize()

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
			// update the nearest intersection
			nearest = distance
			nearestIntersection = &Intersection{
				Point:    ray.Origin.Add(ray.Direction.MultiplyScalar(distance)),
				Distance: distance,
				Object:   object,
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

func Shade(scene *Scene, ray *Ray, intersection *Intersection) *Color {
	//calculate the surface normal at the intersection point
	normal := intersection.Object.Normal(&intersection.Point)

	// initialize the pixel color to the ambient color of the material
	color := intersection.Object.Material().Color.MultiplyScalar(scene.AmbientLight)

	// iterate over all the lights in the scene
	for _, light := range scene.Lights {
		// calculate the direction from the intersection point to the light source
		lightDir := light.Position.Subtract(intersection.Point).Normalize()

		// calculate the diffuse and specular components of the lighting
		diffuse := normal.Dot(lightDir)
		if diffuse > 0 {
			// calculate the specular component
			reflectDir := lightDir.Reflect(normal)
			specular := ray.Direction.Dot(reflectDir)
			if specular > 0 {
				specular = math.Pow(specular, intersection.Object.Material().Shininess)
			} else {
				specular = 0
			}

			// add the diffuse and specular components to the pixel color
			color = color.Add(intersection.Object.Material().Color.MultiplyScalar(diffuse).MultiplyScalar(light.Intensity))
			color = color.Add(light.Color.MultiplyScalar(intersection.Object.Material().Specular).MultiplyScalar(specular).MultiplyScalar(light.Intensity))
		}
	}

	color.R = Clamp(color.R, 0, 1)
	color.G = Clamp(color.G, 0, 1)
	color.B = Clamp(color.B, 0, 1)

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
