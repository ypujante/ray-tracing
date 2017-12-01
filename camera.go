package main

import (
	"math"
)

type Camera interface {
	ray(rnd Rnd, u, v float64) *Ray
}

type camera struct {
	origin          Point3
	lowerLeftCorner Point3
	horizontal      Vec3
	vertical        Vec3
	u, v            Vec3
	lensRadius      float64
}

func NewCamera(lookFrom Point3, lookAt Point3, vup Vec3, vfov float64, aspect float64, aperture float64, focusDist float64) Camera {
	theta := vfov * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := aspect * halfHeight

	origin := lookFrom
	w := lookFrom.Sub(lookAt).Unit()
	u := Cross(vup, w).Unit()
	v := Cross(w, u)

	lowerLeftCorner := origin.Translate(u.Scale(-(halfWidth * focusDist))).Translate(v.Scale(-(halfHeight * focusDist))).Translate(w.Scale(-focusDist))
	horizontal := u.Scale(2 * halfWidth * focusDist)
	vertical := v.Scale(2 * halfHeight * focusDist)

	return camera{origin, lowerLeftCorner, horizontal, vertical, u, v, aperture / 2.0}
}

func (c camera) ray(rnd Rnd, u, v float64) *Ray {
	d := c.lowerLeftCorner.Translate(c.horizontal.Scale(u)).Translate(c.vertical.Scale(v)).Sub(c.origin)
	origin := c.origin

	if c.lensRadius > 0 {
		rd := randomInUnitDisk(rnd).Scale(c.lensRadius)
		offset := c.u.Scale(rd.X).Add(c.v.Scale(rd.Y))
		origin = origin.Translate(offset)
		d = d.Sub(offset)
	}
	return &Ray{origin, d, rnd}
}
