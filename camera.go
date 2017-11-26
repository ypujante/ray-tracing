package main

type Camera struct {
	origin Point3
	lowerLeftCorner Point3
	horizontal Vec3
	vertical Vec3
}

func (c Camera) ray(u, v float64) *Ray {
	d := c.lowerLeftCorner.Translate(c.horizontal.Scale(u)).Translate(c.vertical.Scale(v)).Sub(c.origin)
	return &Ray{c.origin, d}
}
