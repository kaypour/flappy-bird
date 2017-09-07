package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

func main() {
	const (
		windowTitle = "Flappy bird!"
	)

	defer func() {
		// handle panic so it closes nicely, rather than leaving a stack trace
		if x := recover(); x != nil {
			fmt.Fprintf(os.Stderr, "Unexpected failure: %v\n", x)
			os.Exit(2)
		}
	}()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(fmt.Errorf("could not initialize sdl: %v", err))
	}
	defer sdl.Quit()

	if err := mix.Init(mix.INIT_OGG); err != nil {
		panic(fmt.Errorf("could not initialize mix: %v", err))
	}
	defer mix.Quit()
	defer mix.CloseAudio()

	if err := mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		panic(fmt.Errorf("could not open audio: %v", err))
	}

	if err := ttf.Init(); err != nil {
		panic(fmt.Errorf("could not initialize ttf: %v", err))
	}
	defer ttf.Quit()

	window, err := sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(fmt.Errorf("could not create sdl window: %v", err))
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(fmt.Errorf("could not create sdl renderer: %v", err))
	}
	defer renderer.Destroy()

	scene, err := newScene(renderer, "data")
	if err != nil {
		panic(fmt.Errorf("could not create scene: %v", err))
	}
	defer scene.destroy()

	if err := scene.run(); err != nil {
		panic(fmt.Errorf("game engine failed: %v", err))
	}
}
