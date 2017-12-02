# ray-tracing

This is a go implementation of the [Ray Tracing in One Weekend](http://amzn.to/2kiJjn7) book.

![Ray Tracing in Action](./images/ray-tracing.gif)

## Motivation

After following the book in C++ and wanting to learn go, I decided to use it at my first real hands on project with go. I wanted to make the code parallel and display the progress.

## Enhancements

* displays the image as it is being rendered (uses SDL)
* processes the image in multiple goroutines and multiple passes (for example, a first pass with 1 ray per pixel so that the rendering happens very quickly, and then further passes with more rays per pixel to enhance the result)
* choose the seed so that the end result is reproducible

## Installation

`go get -v github.com/ypujante/ray-tracing`

The only dependency is on the [SDL2 Binding for Go](https://github.com/veandco/go-sdl2) library which also requires that you [install SDL](https://github.com/veandco/go-sdl2#requirements) for your system. 

## Compiling

`go install`

Love the simplicity...

## Running

* `ray-tracing` will render the final scene in 800x400 pixels in 2 passes (a quick 1 ray per pixel pass and a 99 rays per pixel pass) using all the cores on your machine and a seed of 2017

* `ray-tracing -r 1 -r 10 -r 50 -r 100 -w 1600 -h 800 -cpu 4 -seed 12345` will use 4 passes (1/10/50/100 rays each so a total of 161 rays per pixel) using `4` cores and a width/height of `1600x800` and a seed of `12345`

## Lessons learned

* `rand.Float64()` is (in hindsight for obvious reasons) synchronized and really killed the performances of the program since it is heavily used by each computation. Abstracted it into a `Rnd` interface (see [model.go](./model.go)) and each goroutine creates its own [non synchronized version](./scene.go#L132) to fix the issue.

* using goroutines and channels really rocks. A few very powerful set of primitives are all it takes to make asynchronous programming a joy again :).

* goroutines are fast... try `-cpu 200`... which will create 200 goroutines and, although you obviously do not gain by having more goroutines than the number of cores (since it is 100% CPU bound), it is pretty fascinating to also see that it is not dramatically slowing it down by having to time slice all those goroutines

* go implements interface/object orientation in a very different manner than any other language. Although it takes some time to get used to it, I really enjoyed it after a while. My `Rnd` interface is a good example, since I could make the `rnd.Rand` class _magically_ implement it even if it is a type not defined by me.

* I do miss generics :( As far as I can tell there is no way to implement the [`split`](./scene.go#L35) function I wrote in a generic fashion which is a shame.

## Dependencies

* [SDL2 Binding for Go](https://github.com/veandco/go-sdl2)

## License

Apache 2.0