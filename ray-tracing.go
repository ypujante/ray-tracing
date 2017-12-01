// Simple ray tracer based on the Ray Tracing book series by Peter Shirley (Kindle)
package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"unsafe"
	"math"
	"runtime"
	"fmt"
	"math/rand"
)

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

func color(r *Ray, world Hitable, depth int) Color {

	if hit, hr := world.hit(r, 0.001, math.MaxFloat64); hit {
		if depth >= 50 {
			return Black
		}

		if wasScattered, attenuation, scattered := hr.material.scatter(r, hr); wasScattered {
			return attenuation.Mult(color(scattered, world, depth+1))
		} else {
			return Black
		}
	}

	unitDirection := r.Direction.Unit()
	t := 0.5 * (unitDirection.Y + 1.0)

	return White.Scale(1.0 - t).Add(Color{0.5, 0.7, 1.0}.Scale(t))
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
			material: Lambertian{Color{0.4, 0.2,0.1}}},
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

func main() {
	const WIDTH, HEIGHT = 800, 400
	//RAYS_PER_PIXEL := []int{2, 4, 4, 10, 10, 10, 10, 10, 10, 10, 10, 10}

	RAYS_PER_PIXEL := []int{1, 99}

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

	camera, world := buildWorldOneWeekend(WIDTH, HEIGHT)
	scene := &Scene{width: WIDTH, height: HEIGHT, raysPerPixel: RAYS_PER_PIXEL, camera: camera, world: world}
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
			break
		default:
			break
		}

		if updateDisplay {
			display(window, screen, scene, pixels)
		}
	}
}
