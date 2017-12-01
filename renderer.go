package main

import (
	"math"
	"math/rand"
	"sync"
	"fmt"
	"time"
)

type PixelToProcess struct {
	i, j, k      int
	color        Color
	raysPerPixel int
}

func split(buf []*PixelToProcess, lim int) [][]*PixelToProcess {
	var chunk []*PixelToProcess
	chunks := make([][]*PixelToProcess, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
	return chunks
}

func (scene *Scene) render(rnd Rnd, pixel *PixelToProcess, raysPerPixel int) uint32 {
	c := pixel.color

	for s := 0; s < raysPerPixel; s++ {
		u := (float64(pixel.i) + rnd.Float64()) / float64(scene.width)
		v := (float64(pixel.j) + rnd.Float64()) / float64(scene.height)
		r := scene.camera.ray(rnd, u, v)
		c = c.Add(color(r, scene.world, 0))
	}

	pixel.color = c
	pixel.raysPerPixel += raysPerPixel

	c = c.Scale(1.0 / float64(pixel.raysPerPixel))
	c = Color{R: math.Sqrt(c.R), G: math.Sqrt(c.G), B: math.Sqrt(c.B)}

	return c.PixelValue()
}

func render(scene *Scene, parallelCount int) (Pixels, chan struct{}) {
	pixels := make([]uint32, scene.width*scene.height)
	completed := make(chan struct{})

	go func() {
		allPixelsToProcess := make([]*PixelToProcess, scene.width*scene.height)

		k := 0
		for j := scene.height - 1; j >= 0; j-- {
			for i := 0; i < scene.width; i++ {
				allPixelsToProcess[k] = &PixelToProcess{i: i, j: j, k: k}
				k++
			}
		}

		lines := split(allPixelsToProcess, scene.width)

		totalRaysPerPixel := 0

		for _, rpp := range scene.raysPerPixel {
			totalRaysPerPixel += rpp
		}

		totalStart := time.Now()

		accumulatedRaysPerPixel := 0

		for _, rpp := range scene.raysPerPixel {

			loopStart := time.Now()

			pixelsToProcess := make(chan []*PixelToProcess)

			go func() {
				for _, p := range lines {
					pixelsToProcess <- p
				}

				close(pixelsToProcess)
			}()

			wg := sync.WaitGroup{}
			for c := 0; c < parallelCount; c++ {
				wg.Add(1)
				go func() {
					rnd := rand.New(rand.NewSource(rand.Int63()))
					for ps := range pixelsToProcess {

						// redisplay the line without gamma correction => make it darker to be more visible
						for _, p := range ps {
							if p.raysPerPixel > 0 {
								col := p.color.Scale(1.0 / float64(p.raysPerPixel))
								pixels[p.k] = col.PixelValue()
							}
						}

						// render every pixel in the line
						for _, p := range ps {
							pixels[p.k] = scene.render(rnd, p, rpp)
						}
					}
					wg.Done()
				}()
			}

			wg.Wait()

			accumulatedRaysPerPixel += rpp

			loopEnd := time.Now()

			totalTimeSoFar := loopEnd.Sub(totalStart)
			estimatedTotalTime := time.Duration(float64(totalTimeSoFar) * float64(totalRaysPerPixel) / float64(accumulatedRaysPerPixel))
			erm := estimatedTotalTime - totalTimeSoFar

			fmt.Printf("Processed %v rays per pixel in %v. Total %v in %v. ERM %v\n", rpp, time.Now().Sub(loopStart), accumulatedRaysPerPixel, totalTimeSoFar, erm)
		}

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
