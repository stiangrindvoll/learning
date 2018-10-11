package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var winTitle = "Guess my pick"
var winWidth, winHeight int32 = 800, 600

var mouseX, mouseY int32
var mouseState uint32
var mousePoint sdl.Point

var selectedCube *sdl.Rect

var fields = createSquares(50, 50, 150, 150, 10)


func run() int {
	var window *sdl.Window
	var renderer *sdl.Renderer

	window, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer renderer.Destroy()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Quit")
				return 0
			}
		}

		mouseX, mouseY, mouseState = sdl.GetMouseState()
		mousePoint = sdl.Point{X: mouseX, Y: mouseY}
		// fmt.Println("mouseX:", mouseX, "mouseY:", mouseY, "mouseState:", mouseState)

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		if mouseState == 1 {
			fields.setSelected(mousePoint)
		}

		fields.render(renderer, mousePoint)

		renderer.Present()

		sdl.Delay(50)

	}
}

func main() {
	os.Exit(run())
}
