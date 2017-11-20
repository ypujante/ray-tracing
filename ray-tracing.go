// Simple ray tracer based on the Ray Tracing book series by Peter Shirley (Kindle)
package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"unsafe"
)

type RenderBlock struct {
	X, Y, Width, Height int
	pixels              []uint32
}

func hitSphere(center Point3, radius float64, ray Ray) bool {
	oc := ray.Origin.Sub(center)
	a := Dot(ray.Direction, ray.Direction)
	b := 2.0 * Dot(oc, ray.Direction)
	c := Dot(oc, oc) - radius * radius
	discriminant := b * b - 4.0 * a * c
	return discriminant > 0
}

func render(Width, Height int) RenderBlock {
	pixels := make([]uint32, Width*Height)

	origin := Point3{}
	lowerLeftCorner := Point3{-2.0, -1.0, -1.0}
	horizontal := Vec3{X: 4.0}
	vertical := Vec3{Y: 2.0}

	k := 0
	for j := Height - 1; j >= 0; j-- {
		for i := 0; i < Width; i++ {

			u := float64(i) / float64(Width)
			v := float64(j) / float64(Height)

			d := lowerLeftCorner.Translate(horizontal.Scale(u)).Translate(vertical.Scale(v)).Vec3()

			c := color(Ray{origin, d})

			pixels[k] = c.PixelValue()

			k++
		}
	}

	return RenderBlock{0, 0, Width, Height, pixels}
}

func color(r Ray) Color {
	if hitSphere(Point3{Z: -1.0}, 0.5, r) {
		return Color{R: 1.0}
	}
	unitDirection := r.Direction.Unit()

	t := 0.5 * (unitDirection.Y + 1.0)

	return White.Scale(1.0-t).Add(Color{0.5, 0.7, 1.0}.Scale(t))
}

func main() {
	const WIDTH, HEIGHT = 200, 100

	// initializes SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// create (and show) window
	window, err := sdl.CreateWindow("Ray Tracing", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WIDTH, HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// retrieves the screen
	screen, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	// clear the screen (otherwise there is garbage...)
	err = screen.FillRect(&sdl.Rect{W: WIDTH, H: HEIGHT}, 0x00000000)
	if err != nil {
		panic(err)
	}

	// actual work to render the image
	rb := render(WIDTH, HEIGHT)

	// create an image from the pixels generated
	image, err := sdl.CreateRGBSurfaceFrom(unsafe.Pointer(&rb.pixels[0]), int32(rb.Width), int32(rb.Height), 32, rb.Width*int(unsafe.Sizeof(rb.pixels[0])), 0, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	// copy it into the screen
	err = image.Blit(nil, screen, nil)
	if err != nil {
		panic(err)
	}
	image.Free()

	// update the surface to show it
	err = window.UpdateSurface()
	if err != nil {
		panic(err)
	}

	// poll for quit event
	for running := true; running; {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		sdl.Delay(16)
	}
}
