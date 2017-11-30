// Simple ray tracer based on the Ray Tracing book series by Peter Shirley (Kindle)
package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"unsafe"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"fmt"
	"time"
)

type Pixels []uint32

type Scene struct {
	width, height int
	raysPerPixel  int
	camera        Camera
	world         Hitable
}

func randomInUnitSphere(rnd *rand.Rand) Vec3 {
	for {
		p := Vec3{2.0*rnd.Float64() - 1.0, 2.0*rnd.Float64() - 1.0, 2.0*rnd.Float64() - 1.0}
		if Dot(p, p) < 1.0 {
			return p
		}
	}
}

func randomInUnitDisk(rnd *rand.Rand) Vec3 {
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

func (scene *Scene) render(rnd *rand.Rand, i, j int) uint32 {
	c := Color{}

	for s := 0; s < scene.raysPerPixel; s++ {
		u := (float64(i) + rnd.Float64()) / float64(scene.width)
		v := (float64(j) + rnd.Float64()) / float64(scene.height)
		r := scene.camera.ray(rnd, u, v)
		c = c.Add(color(r, scene.world, 0))
	}

	c = c.Scale(1.0 / float64(scene.raysPerPixel))
	c = Color{R: math.Sqrt(c.R), G: math.Sqrt(c.G), B: math.Sqrt(c.B)}

	return c.PixelValue()
}


func render(scene *Scene, parallelCount int) (Pixels, chan struct{}) {
	pixels := make([]uint32, scene.width*scene.height)

	completed := make(chan struct{})

	type PixelToProcess struct {
		i,j,k int
	}

	split := func (buf []PixelToProcess, lim int) [][]PixelToProcess {
		var chunk []PixelToProcess
		chunks := make([][]PixelToProcess, 0, len(buf)/lim+1)
		for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
		if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
		return chunks
	}


	pixelsToProcess := make(chan []PixelToProcess)


	allPixelsToProcess := make([]PixelToProcess, scene.width*scene.height)

	k := 0
	for j := scene.height - 1; j >= 0; j-- {
		for i := 0; i < scene.width; i++ {
			allPixelsToProcess[k] = PixelToProcess{i,j,k}
			k++
		}
	}

	go func() {
		slices := split(allPixelsToProcess, scene.width)

		for _, p := range slices {
			pixelsToProcess <- p
		}

		close(pixelsToProcess)
	}()

	wg := sync.WaitGroup{}

	for c:= 0; c < parallelCount; c++ {
		wg.Add(1)
		go func() {
			rnd := rand.New(rand.NewSource(1971))
			for ps := range pixelsToProcess {
				for _, p:= range ps {
					pixels[p.k] = scene.render(rnd, p.i, p.j)
				}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		completed <- struct{}{}
	}()

	return pixels, completed
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

func buildWorldDielectrics() HitableList {
	return HitableList{
		Sphere{center: Point3{Z: -1.0}, radius: 0.5, material: Lambertian{Color{R: 0.1, G: 0.2, B: 0.5}}},
		Sphere{center: Point3{Y: -100.5, Z: -1.0}, radius: 100, material: Lambertian{Color{R: 0.8, G: 0.8}}},
		Sphere{center: Point3{X: 1.0, Y: 0, Z: -1.0}, radius: 0.5, material: Metal{Color{R: 0.8, G: 0.6, B: 0.2}, 1.0}},
		Sphere{center: Point3{X: -1.0, Y: 0, Z: -1.0}, radius: 0.5, material: Dielectric{1.5}},
		Sphere{center: Point3{X: -1.0, Y: 0, Z: -1.0}, radius: -0.45, material: Dielectric{1.5}},
	}
}

func main() {
	const WIDTH, HEIGHT, RAYS_PER_PIXEL = 800, 400, 100

	rand.Seed(1971)

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
	lookFrom := Point3{-2.0, 2.0, 1.0}
	lookAt := Point3{Z: -1.0}
	//aperture := 2.0
	//distToFocus := lookFrom.Sub(lookAt).Length()
	aperture := 0.0
	distToFocus := 1.0
	camera := NewCamera(lookFrom, lookAt, Vec3{Y: 1.0}, 20, WIDTH/HEIGHT, aperture, distToFocus)

	world := buildWorldDielectrics()
	scene := &Scene{width: WIDTH, height: HEIGHT, raysPerPixel: RAYS_PER_PIXEL, camera: camera, world: world}
	pixels, completed := render(scene, runtime.NumCPU())

	// update the surface to show it
	err = window.UpdateSurface()
	if err != nil {
		panic(err)
	}

	updateDisplay := true

	start := time.Now()

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
			fmt.Println(time.Now().Sub(start))
			break
		default:
			break
		}

		if updateDisplay {
			display(window, screen, scene, pixels)
		}
	}
}
