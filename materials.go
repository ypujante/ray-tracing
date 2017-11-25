package main

/***********************
 * Material
 ************************/
// Material defines how a material scatter light
type Material interface {
	scatter(r Ray, rec *HitRecord) (wasScattered bool , attenuation *Color, scattered *Ray)
}

/***********************
 * LambertianMaterial material (diffuse only)
 ************************/
type LambertianMaterial struct {
	albedo Color
}

func (l LambertianMaterial) scatter(r Ray, rec *HitRecord) (bool, *Color, *Ray) {
	target := rec.p.Translate(rec.normal).Translate(randomInUnitSphere())
	scattered := &Ray{rec.p, target.Sub(rec.p)}
	attenuation := &l.albedo
	return true, attenuation, scattered

}