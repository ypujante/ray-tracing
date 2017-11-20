// Simple ray tracer based on the Ray Tracing book series by Peter Shirley (Kindle)
package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"unsafe"
)

type RenderBlock struct {
	X, Y, Width, Height int
	pixels              []uint32
}

func render(Width, Height int) RenderBlock {
	pixels := make([]uint32, Width*Height)
	k := 0
	for j := Height - 1; j >= 0; j-- {
		for i := 0; i < Width; i++ {
			r := float64(i) / float64(Width)
			g := float64(j) / float64(Height)
			b := float64(0.2)

			R := uint32(math.Min(255.0, r*255.99))
			G := uint32(math.Min(255.0, g*255.99))
			B := uint32(math.Min(255.0, b*255.99))

			pixels[k] = ((R & 0xFF) << 16) | ((G & 0xFF) << 8) | (B & 0xFF)
			k++
		}
	}

	return RenderBlock{0, 0, Width, Height, pixels}
}

func main() {
	const WIDTH, HEIGHT = 800, 600

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
	err = screen.FillRect(&sdl.Rect{0, 0, WIDTH, HEIGHT}, 0x00000000)
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
