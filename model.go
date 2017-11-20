package main

import "math"

// Vec3 defines a vector in 3D space
type Vec3 struct {
	X, Y, Z float64
}

// Scale scales the vector by the value (return a new vector)
func (v Vec3) Scale(t float64) Vec3 {
	return Vec3{X: v.X * t, Y: v.Y * t, Z: v.Z * t}
}

// Length returns the size of the vector
func (v Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Unit returns a new vector with same direction and length 1
func (v Vec3) Unit() Vec3 {
	return v.Scale(1.0 / v.Length())
}

// Point3 defines a point in 3D space
type Point3 struct {
	X, Y, Z float64
}

// Translate translates the point to a new location (return a new point)
func (p Point3) Translate(v Vec3) Point3 {
	return Point3{p.X + v.X, p.Y + v.Y, p.Z + v.Z}
}

// Vec3 converts a point to a vector (centered at origin)
func (p Point3) Vec3() Vec3 {
	return Vec3{p.X, p.Y, p.Z }
}

// Ray represents a ray defined by its origin and direction
type Ray struct {
	Origin    Point3
	Direction Vec3
}

// PointAt returns a new point along the ray (0 will return the origin)
func (r Ray) PointAt(t float64) Point3 {
	return r.Origin.Translate(r.Direction.Scale(t))
}

// Color defines the basic Red/Green/Blue as raw float64 values
type Color struct {
	R, G, B float64
}

var (
	White = Color{1.0, 1.0, 1.0}
)

// Scale scales the Color by the value (return a new Color)
func (c Color) Scale(t float64) Color {
	return Color{R: c.R * t, G: c.G * t, B: c.B * t}
}

// Add adds the 2 colors (return a new color)
func (c Color) Add(c2 Color) Color {
	return Color{R: c.R + c2.R, G: c.G + c2.G, B: c.B + c2.B}
}

// PixelValue converts a raw Color into a pixel value (0-255) packed into a uint32
func (c Color) PixelValue() uint32 {
	r := uint32(math.Min(255.0, c.R*255.99))
	g := uint32(math.Min(255.0, c.G*255.99))
	b := uint32(math.Min(255.0, c.B*255.99))

	return ((r & 0xFF) << 16) | ((g & 0xFF) << 8) | (b & 0xFF)
}
