package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const winW int = 800
const winH int = 600

type color struct {
	r byte
	g byte
	b byte
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winW + x) * 4

	if index < len(pixels)-4 && index >= 0 {

		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b

	}
}

func main() {
	fmt.Println("vim-go")
	window, err := sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winW), int32(winH), sdl.WINDOW_SHOWN)
	check(err)

	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	check(err)
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winW), int32(winH))
	check(err)
	defer tex.Destroy()

	pixels := make([]byte, winW*winH*4)

	for y := 0; y < winH; y++ {
		for x := 0; x < winW; x++ {
			setPixel(x, y, color{byte(x % 255), byte(x % 255), byte((int(x+y) / 2) % 255)}, pixels)
		}
	}

	tex.Update(nil, pixels, winW*4)
	renderer.Copy(tex, nil, nil)
	renderer.Present()

	sdl.Delay(2000)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
