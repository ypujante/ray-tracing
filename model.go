package main

import (
	"math"
	"math/rand"
)

/***********************
 * Vec3
 ************************/
// Vec3 defines a vector in 3D space
type Vec3 struct {
	X, Y, Z float64
}

// Scale scales the vector by the value (return a new vector)
func (v Vec3) Scale(t float64) Vec3 {
	return Vec3{X: v.X * t, Y: v.Y * t, Z: v.Z * t}
}

// Mult multiplies the vector by the other one (return a new vector)
func (v Vec3) Mult(v2 Vec3) Vec3 {
	return Vec3{X: v.X * v2.X, Y: v.Y * v2.Y, Z: v.Z * v2.Z}
}

// Sub substracts the 2 vectors (return a new vector)
func (v Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{X: v.X - v2.X, Y: v.Y - v2.Y, Z: v.Z - v2.Z}
}

// Add adds the 2 vectors (return a new vector)
func (v Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{X: v.X + v2.X, Y: v.Y + v2.Y, Z: v.Z + v2.Z}
}

// Length returns the size of the vector
func (v Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Unit returns a new vector with same direction and length 1
func (v Vec3) Unit() Vec3 {
	return v.Scale(1.0 / v.Length())
}

// Negate returns a new vector with X/Y/Z negated
func (v Vec3) Negate() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}

// Dot returns the dot product (a scalar) of 2 vectors
func Dot(v1 Vec3, v2 Vec3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

// Cross returns the cross product of 2 vectors (another vector)
func Cross(v1 Vec3, v2 Vec3) Vec3 {
	return Vec3{v1.Y*v2.Z - v1.Z*v2.Y, -(v1.X*v2.Z - v1.Z*v2.X), v1.X*v2.Y - v1.Y*v2.X}
}

// Reflect simply reflects the vector based on the normal n
func (v Vec3) Reflect(n Vec3) Vec3 {
	return v.Sub(n.Scale(2.0 * Dot(v, n)))
}

/***********************
 * Point3
 ************************/
// Point3 defines a point in 3D space
type Point3 struct {
	X, Y, Z float64
}

// Translate translates the point to a new location (return a new point)
func (p Point3) Translate(v Vec3) Point3 {
	return Point3{p.X + v.X, p.Y + v.Y, p.Z + v.Z}
}

// Sub subtracts a point to another Point which gives a vector
func (p Point3) Sub(p2 Point3) Vec3 {
	return Vec3{p.X - p2.X, p.Y - p2.Y, p.Z - p2.Z}
}

// Vec3 converts a point to a vector (centered at origin)
func (p Point3) Vec3() Vec3 {
	return Vec3{p.X, p.Y, p.Z}
}

/***********************
 * Ray
 ************************/
// Ray represents a ray defined by its origin and direction
type Ray struct {
	Origin    Point3
	Direction Vec3
	rnd *rand.Rand
}

// PointAt returns a new point along the ray (0 will return the origin)
func (r *Ray) PointAt(t float64) Point3 {
	return r.Origin.Translate(r.Direction.Scale(t))
}

/***********************
 * Color
 ************************/
// Color defines the basic Red/Green/Blue as raw float64 values
type Color struct {
	R, G, B float64
}

var (
	White = Color{1.0, 1.0, 1.0}
	Black = Color{}
)

// Scale scales the Color by the value (return a new Color)
func (c Color) Scale(t float64) Color {
	return Color{R: c.R * t, G: c.G * t, B: c.B * t}
}

// Mult Multiplies 2 colors together (component by component multiplication)
func (c Color) Mult(c2 Color) Color {
	return Color{R: c.R * c2.R, G: c.G * c2.G, B: c.B * c2.B}
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

/***********************
 * Hitable
 ************************/
type HitRecord struct {
	t        float64  // which t generated the hit
	p        Point3   // which point when hit
	normal   Vec3     // normal at that point
	material Material // the material associated to this record
}

// Hitable defines the interface of objects that can be hit by a ray
type Hitable interface {
	hit(r *Ray, tMin float64, tMax float64) (bool, *HitRecord)
}

// HitableList defines a simple list of hitable
type HitableList []Hitable

// hit defines the method for a list of hitables: will return the one closest
func (hl HitableList) hit(r *Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	var res *HitRecord
	hitAnything := false

	closestSoFar := tMax

	for _, h := range hl {
		if hit, hr := h.hit(r, tMin, closestSoFar); hit {
			hitAnything = true
			res = hr
			closestSoFar = hr.t
		}
	}

	return hitAnything, res
}
