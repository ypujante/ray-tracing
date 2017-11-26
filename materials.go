package main

import (
	"math"
	"math/rand"
)

/***********************
 * Material
 ************************/
// Material defines how a material scatter light
type Material interface {
	scatter(r *Ray, rec *HitRecord) (wasScattered bool, attenuation *Color, scattered *Ray)
}

/***********************
 * Lambertian material (diffuse only)
 ************************/
type Lambertian struct {
	albedo Color
}

func (mat Lambertian) scatter(r *Ray, rec *HitRecord) (bool, *Color, *Ray) {
	target := rec.p.Translate(rec.normal).Translate(randomInUnitSphere())
	scattered := &Ray{rec.p, target.Sub(rec.p)}
	attenuation := &mat.albedo
	return true, attenuation, scattered

}

/***********************
 * Metal material
 ************************/
type Metal struct {
	albedo Color
	fuzz   float64
}

func (mat Metal) scatter(r *Ray, rec *HitRecord) (bool, *Color, *Ray) {
	reflected := r.Direction.Unit().Reflect(rec.normal)
	if mat.fuzz < 1 {
		reflected = reflected.Add(randomInUnitSphere().Scale(mat.fuzz))
	}
	scattered := &Ray{rec.p, reflected}
	attenuation := &mat.albedo

	if Dot(scattered.Direction, rec.normal) > 0 {
		return true, attenuation, scattered
	}

	return false, nil, nil
}

/***********************
 * Dielectric material (glass)
 ************************/
type Dielectric struct {
	refIdx float64
}

// Refract returns a refracted vector (or not if there is no refraction possible)
func (v Vec3) Refract(n Vec3, niOverNt float64) (bool, *Vec3) {
	uv := v.Unit()
	un := n.Unit()

	dt := Dot(uv, un)
	discriminant := 1.0 - niOverNt*niOverNt*(1-dt*dt)
	if discriminant > 0 {
		refracted := uv.Sub(un.Scale(dt)).Scale(niOverNt).Sub(un.Scale(math.Sqrt(discriminant)))
		return true, &refracted
	}

	return false, nil
}

func schlick(cosine float64, iRefIdx float64) float64 {
	r0 := (1.0 - iRefIdx) / (1.0 + iRefIdx)
	r0 = r0 * r0
	return r0 + (1.0-r0)*math.Pow(1.0-cosine, 5)
}

func (die Dielectric) scatter(r *Ray, rec *HitRecord) (bool, *Color, *Ray) {
	var (
		outwardNormal Vec3
		niOverNt      float64
		cosine        float64
	)

	dotRayNormal := Dot(r.Direction, rec.normal);
	if dotRayNormal > 0 {
		outwardNormal = rec.normal.Negate()
		niOverNt = die.refIdx;
		cosine = dotRayNormal / r.Direction.Length()
		cosine = math.Sqrt(1.0 - die.refIdx*die.refIdx*(1.0-cosine*cosine));
	} else {
		outwardNormal = rec.normal;
		niOverNt = 1.0 / die.refIdx
		cosine = -dotRayNormal / r.Direction.Length()
	}

	wasRefracted, refracted := r.Direction.Refract(outwardNormal, niOverNt)

	var direction Vec3

	// refract only with some probability
	if wasRefracted && rand.Float64() >= schlick(cosine, die.refIdx) {
		direction = *refracted
	} else {
		direction = r.Direction.Unit().Reflect(rec.normal)
	}

	return true, &White, &Ray{rec.p, direction}
}
