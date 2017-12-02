package main

import "math"

type Sphere struct {
	center   Point3
	radius   float64
	material Material
}

// hit implements the hit interface for a Sphere
func (s Sphere) hit(r *Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	oc := r.Origin.Sub(s.center)         // A-C
	a := Dot(r.Direction, r.Direction)   // dot(B, B)
	b := Dot(oc, r.Direction)            // dot(A-C, B)
	c := Dot(oc, oc) - s.radius*s.radius // dot(A-C, A-C) - R*R
	discriminant := b*b - a*c

	if discriminant > 0 {
		discriminantSquareRoot := math.Sqrt(discriminant)

		temp := (-b - discriminantSquareRoot) / a
		if temp < tMax && temp > tMin {
			hitPoint := r.PointAt(temp)
			hr := HitRecord{
				t:        temp,
				p:        hitPoint,
				normal:   hitPoint.Sub(s.center).Scale(1 / s.radius),
				material: s.material,
			}
			return true, &hr
		}

		temp = (-b + discriminantSquareRoot) / a
		if temp < tMax && temp > tMin {
			hitPoint := r.PointAt(temp)
			hr := HitRecord{
				t:        temp,
				p:        hitPoint,
				normal:   hitPoint.Sub(s.center).Scale(1 / s.radius),
				material: s.material,
			}
			return true, &hr
		}
	}

	return false, nil

}
