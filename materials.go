package main

/***********************
 * Material
 ************************/
// Material defines how a material scatter light
type Material interface {
	scatter(r *Ray, rec *HitRecord) (wasScattered bool , attenuation *Color, scattered *Ray)
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
 * MetalMaterial material
 ************************/
type Metal struct {
	albedo Color
}

func (mat Metal) scatter(r *Ray, rec *HitRecord) (bool, *Color, *Ray) {
	reflected := r.Direction.Unit().Reflect(rec.normal)
	scattered := &Ray{rec.p, reflected}
	attenuation := &mat.albedo

	if Dot(scattered.Direction, rec.normal) > 0 {
		return true, attenuation, scattered
	}

	return false, nil, nil
}

