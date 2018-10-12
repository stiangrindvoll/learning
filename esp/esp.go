package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var winTitle = "Guess my pick"
var winWidth, winHeight int32 = 800, 600

var mouseX, mouseY int32
var mouseState uint32
var mousePoint sdl.Point

type player struct {
	name  string
	field field
}

var selectedCube *sdl.Rect

var communication *bufio.ReadWriter

func run(name string) int {
	var window *sdl.Window
	var renderer *sdl.Renderer

	player := player{name: name, field: field{createSquares(50, 50, 150, 150, 10), -1}}

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
			player.field.setSelected(mousePoint, communication)
		}

		player.field.render(renderer, mousePoint)

		renderer.Present()

		sdl.Delay(50)

	}
}

func main() {
	help := flag.Bool("h", false, "Display Help")
	rendezvousString := flag.String("r", "ESPGAME", "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	name := flag.String("n", "ESP Game", "Player name")
	flag.Parse()

	if *help {
		fmt.Printf("This program demonstrates a simple p2p chat application using libp2p\n\n")
		fmt.Printf("Usage: Run './chat in two different terminals. Let them connect to the bootstrap nodes, announce themselves and connect to the peers\n")

		os.Exit(0)
	}
	go createNetwork(*rendezvousString, communication)
	os.Exit(run(*name))
}
