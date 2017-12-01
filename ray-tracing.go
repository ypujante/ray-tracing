// Simple ray tracer based on the Ray Tracing book series by Peter Shirley (Kindle)
package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"unsafe"
	"runtime"
	"fmt"
	"math/rand"
	"flag"
	"strings"
	"strconv"
	"image"
	clr "image/color"
	"os"
	"image/png"
)

type RaysPerPixelList []int

type Options struct {
	Width        int
	Height       int
	RaysPerPixel RaysPerPixelList
	Output       string
}

type Pixels []uint32

type Scene struct {
	width, height int
	raysPerPixel  []int
	camera        Camera
	world         Hitable
}

func randomInUnitSphere(rnd Rnd) Vec3 {
	for {
		p := Vec3{2.0*rnd.Float64() - 1.0, 2.0*rnd.Float64() - 1.0, 2.0*rnd.Float64() - 1.0}
		if Dot(p, p) < 1.0 {
			return p
		}
	}
}

func randomInUnitDisk(rnd Rnd) Vec3 {
	for {
		p := Vec3{2.0*rnd.Float64() - 1.0, 2.0*rnd.Float64() - 1.0, 0}
		if Dot(p, p) < 1.0 {
			return p
		}
	}
}

func display(window *sdl.Window, screen *sdl.Surface, scene *Scene, pixels Pixels) {
	// create an image from the pixels generated
	image, err := sdl.CreateRGBSurfaceFrom(unsafe.Pointer(&pixels[0]), int32(scene.width), int32(scene.height), 32, scene.width*int(unsafe.Sizeof(pixels[0])), 0, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	defer image.Free()
	// copy it into the screen
	err = image.Blit(nil, screen, nil)
	if err != nil {
		panic(err)
	}

	// update the surface to show it
	err = window.UpdateSurface()
	if err != nil {
		panic(err)
	}
}

func buildWorld() HitableList {
	return HitableList{
		Sphere{center: Point3{Z: -1.0}, radius: 0.5, material: Lambertian{Color{R: 1.0}}},
		Sphere{center: Point3{Y: -100.5, Z: -1.0}, radius: 100, material: Lambertian{Color{G: 1.0}}},
	}
}

func buildWorldMetalSpheres() HitableList {
	return HitableList{
		Sphere{center: Point3{Z: -1.0}, radius: 0.5, material: Lambertian{Color{R: 0.8, G: 0.3, B: 0.3}}},
		Sphere{center: Point3{Y: -100.5, Z: -1.0}, radius: 100, material: Lambertian{Color{R: 0.8, G: 0.8}}},
		Sphere{center: Point3{X: 1.0, Y: 0, Z: -1.0}, radius: 0.5, material: Metal{Color{R: 0.8, G: 0.6, B: 0.2}, 1.0}},
		Sphere{center: Point3{X: -1.0, Y: 0, Z: -1.0}, radius: 0.5, material: Metal{Color{R: 0.8, G: 0.8, B: 0.8}, 0.3}},
	}
}

func buildWorldDielectrics(width, height int) (Camera, HitableList) {

	lookFrom := Point3{-2.0, 2.0, 1.0}
	lookAt := Point3{Z: -1.0}
	aperture := 0.0
	distToFocus := 1.0
	camera := NewCamera(lookFrom, lookAt, Vec3{Y: 1.0}, 20, float64(width)/float64(height), aperture, distToFocus)

	world := HitableList{
		Sphere{center: Point3{Z: -1.0}, radius: 0.5, material: Lambertian{Color{R: 0.1, G: 0.2, B: 0.5}}},
		Sphere{center: Point3{Y: -100.5, Z: -1.0}, radius: 100, material: Lambertian{Color{R: 0.8, G: 0.8}}},
		Sphere{center: Point3{X: 1.0, Y: 0, Z: -1.0}, radius: 0.5, material: Metal{Color{R: 0.8, G: 0.6, B: 0.2}, 1.0}},
		Sphere{center: Point3{X: -1.0, Y: 0, Z: -1.0}, radius: 0.5, material: Dielectric{1.5}},
		Sphere{center: Point3{X: -1.0, Y: 0, Z: -1.0}, radius: -0.45, material: Dielectric{1.5}},
	}

	return camera, world
}

func buildWorldOneWeekend(width, height int) (Camera, HitableList) {
	world := []Hitable{}

	maxSpheres := 500
	world = append(world, Sphere{center: Point3{Y: -1000.0}, radius: 1000, material: Lambertian{Color{R: 0.5, G: 0.5, B: 0.5}}})

	rand.Seed(1971)

	for a := -11; a < 11 && len(world) < maxSpheres; a++ {
		for b := -11; b < 11 && len(world) < maxSpheres; b++ {
			chooseMaterial := rand.Float64()
			center := Point3{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}

			if center.Sub(Point3{4.0, 0.2, 0}).Length() > 0.9 {
				switch {
				case chooseMaterial < 0.8: // diffuse
					world = append(world,
						Sphere{
							center:   center,
							radius:   0.2,
							material: Lambertian{Color{R: rand.Float64() * rand.Float64(), G: rand.Float64() * rand.Float64(), B: rand.Float64() * rand.Float64()}}})
				case chooseMaterial < 0.95: // metal
					world = append(world,
						Sphere{
							center:   center,
							radius:   0.2,
							material: Metal{Color{R: 0.5 * (1 + rand.Float64()), G: 0.5 * (1 + rand.Float64()), B: 0.5 * (1 + rand.Float64())}, 0.5 * rand.Float64()}})
				default:
					world = append(world,
						Sphere{
							center:   center,
							radius:   0.2,
							material: Dielectric{1.5}})

				}
			}
		}
	}

	world = append(world,
		Sphere{
			center:   Point3{0, 1, 0},
			radius:   1.0,
			material: Dielectric{1.5}},
		Sphere{
			center:   Point3{-4, 1, 0},
			radius:   1.0,
			material: Lambertian{Color{0.4, 0.2, 0.1}}},
		Sphere{
			center:   Point3{4, 1, 0},
			radius:   1.0,
			material: Metal{Color{0.7, 0.6, 0.5}, 0}})

	lookFrom := Point3{13, 2, 3}
	lookAt := Point3{}
	aperture := 0.1
	distToFocus := 10.0
	camera := NewCamera(lookFrom, lookAt, Vec3{Y: 1.0}, 20, float64(width)/float64(height), aperture, distToFocus)

	return camera, world
}

func (r *RaysPerPixelList) String() string {
	return fmt.Sprint(*r)
}

func (r *RaysPerPixelList) Set(value string) error {
	for _, e := range strings.Split(value, ",") {
		i, err := strconv.Atoi(e)
		if err != nil {
			return err
		}
		*r = append(*r, i)
	}
	return nil
}

func saveImage(pixels Pixels, options Options) (error, bool) {
	if options.Output != "" {
		f, err := os.OpenFile(options.Output, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err, true
		}

		img := image.NewNRGBA(image.Rect(0, 0, options.Width, options.Height))

		k := 0
		for y := 0; y < options.Height; y++ {
			for x := 0; x < options.Width; x++ {
				p := pixels[k]
				img.Set(x, y, clr.NRGBA{
					R: uint8(p >> 16 & 0xFF),
					G: uint8(p >> 8 & 0xFF),
					B: uint8(p & 0xFF),
					A: 255,
				})
				k++
			}
		}

		if err := png.Encode(f, img); err != nil {
			f.Close()
			return err, true
		}

		if err := f.Close(); err != nil {
			return err, true
		}

		return nil, true
	}

	return nil, false

}

func main() {
	options := Options{}

	flag.IntVar(&options.Width, "w", 800, "width in pixel")
	flag.IntVar(&options.Height, "h", 400, "height in pixel")
	flag.Var(&options.RaysPerPixel, "r", "comma separated list (or multiple) rays per pixel")
	flag.StringVar(&options.Output, "o", "", "path to file for saving (do not save if not defined)")

	flag.Parse()

	if (len(options.RaysPerPixel) == 0) {
		options.RaysPerPixel = []int{1, 99}
	}

	// initializes SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// create (and show) window
	window, err := sdl.CreateWindow("Ray Tracing", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(options.Width), int32(options.Height), sdl.WINDOW_SHOWN)
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
	err = screen.FillRect(&sdl.Rect{W: int32(options.Width), H: int32(options.Height)}, 0x00000000)
	if err != nil {
		panic(err)
	}

	camera, world := buildWorldOneWeekend(options.Width, options.Height)
	scene := &Scene{width: options.Width, height: options.Height, raysPerPixel: options.RaysPerPixel, camera: camera, world: world}
	pixels, completed := render(scene, runtime.NumCPU())

	// update the surface to show it
	err = window.UpdateSurface()
	if err != nil {
		panic(err)
	}

	updateDisplay := true

	// poll for quit event
	for running := true; running; {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		sdl.Delay(16)

		select {
		case <-completed:
			display(window, screen, scene, pixels)
			updateDisplay = false
			fmt.Println("Render complete.")
			err, saved := saveImage(pixels, options)
			switch {
			case err != nil:
				fmt.Printf("Error while saving the image [%v]\n", err)
			case saved:
				fmt.Printf("Image saved to %v\n", options.Output)
			}
		default:
			break
		}

		if updateDisplay {
			display(window, screen, scene, pixels)
		}
	}
}
