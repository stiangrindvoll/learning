package main

import (
	"fmt"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const winW int = 1280
const winH int = 720

type gameState int

const (
	start gameState = iota
	play
)

var state = start

var nums = [][]byte{
	{
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	{
		1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		1, 0, 1,
		1, 0, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		0, 1, 1,
		0, 0, 1,
		1, 1, 1}}

type color struct {
	r byte
	g byte
	b byte
}

type ball struct {
	pos
	radius float32
	xv     float32
	yv     float32
	color  color
}

type key struct {
	up, down sdl.Scancode
}

func drawNumber(pos pos, color color, size int, num int, pixels []byte) {
	startX := int(pos.x) - (size*3)/2
	startY := int(pos.y) - (size*5)/2

	for i, v := range nums[num] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += size
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}

func lerp(a, b, p float32) float32 {
	return a + p*(b-a)

}

func (b *ball) draw(pixels []byte) {
	// draw only simple circle within a radius of a square
	for y := -b.radius; y < b.radius; y++ {
		for x := -b.radius; x < b.radius; x++ {
			if x*x+y*y < b.radius*b.radius {
				setPixel(int(b.x+x), int(b.y+y), b.color, pixels)
			}
		}
	}
}

func (b *ball) update(leftPaddle, rightPaddle *paddle, elapsedTime float32) {
	b.x += b.xv * elapsedTime
	b.y += b.yv * elapsedTime

	// handle top/bottom collisions
	if b.y+-b.radius < 0 {
		b.yv = -b.yv
		b.y = b.radius
	} else if b.y+b.radius > float32(winH) {
		b.yv = -b.yv
		b.y = b.y - b.radius
	}

	// handle paddle collisions
	if b.x-b.radius < leftPaddle.x+leftPaddle.w/2 {
		if b.y > leftPaddle.y-leftPaddle.h/2 && b.y < leftPaddle.y+leftPaddle.h/2 {
			b.xv = -b.xv
			b.x = leftPaddle.x + leftPaddle.w/2.0 + b.radius

		}
	}
	if b.x+b.radius > rightPaddle.x-leftPaddle.w/2 {
		if b.y > rightPaddle.y-rightPaddle.h/2 && b.y < rightPaddle.y+rightPaddle.h/2 {
			b.xv = -b.xv
			b.x = rightPaddle.x - rightPaddle.w/2.0 - b.radius
		}
	}

	if b.x < 0 {
		rightPaddle.score++
		b.pos = getCenter()
		state = start
	} else if int(b.x) > winW {
		leftPaddle.score++
		b.pos = getCenter()
		state = start
	}

}

func getCenter() pos {
	return pos{float32(winW / 2), float32(winH / 2)}
}

type paddle struct {
	pos
	w, h  float32
	color color
	speed float32
	score int
}

func (p *paddle) update(keyState []uint8, key key, controllerAxis int16, elapsedTime float32) {
	if keyState[key.up] != 0 {
		p.y -= p.speed * elapsedTime
	}
	if keyState[key.down] != 0 {
		p.y += p.speed * elapsedTime
	}

	if math.Abs(float64(controllerAxis)) > 1500 {
		pct := float32(controllerAxis) / 32767.0
		p.y += p.speed * pct * elapsedTime
	}
}

func (p *paddle) aiupdate(b *ball, elapsedTime float32) {
	p.y = b.y
}

func (p *paddle) draw(pixels []byte) {
	// pos is the center
	startX := int(p.x - p.w/2)
	startY := int(p.y - p.h/2)

	for y := 0; y < int(p.h); y++ {
		for x := 0; x < int(p.w); x++ {
			setPixel(startX+x, startY+y, p.color, pixels)
		}
	}
	numX := lerp(p.x, getCenter().x, 0.2)
	drawNumber(pos{numX, 50}, color{100, 100, 100}, 10, p.score, pixels)
}

type pos struct {
	x, y float32
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winW + x) * 4

	if index < len(pixels)-4 && index >= 0 {

		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b

	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	check(err)
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Pong", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winW), int32(winH), sdl.WINDOW_SHOWN)
	check(err)

	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	check(err)
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winW), int32(winH))
	check(err)
	defer tex.Destroy()

	pixels := make([]byte, winW*winH*4)

	var controllerHandlers []*sdl.GameController

	for i := 0; i < sdl.NumJoysticks(); i++ {
		controllerHandlers = append(controllerHandlers, sdl.GameControllerOpen(i))
		defer controllerHandlers[i].Close()

	}

	//for y := 0; y < winH; y++ {
	//	for x := 0; x < winW; x++ {
	//		setPixel(x, y, color{byte(x % 255), byte(x % 255), byte((int(x+y) / 2) % 255)}, pixels)
	//	}
	//}

	player1 := paddle{pos{100, 100}, 20, 100, color{255, 255, 255}, 600, 0}
	player2 := paddle{pos{float32(winW) - 100, 100}, 20, 100, color{255, 255, 255}, 600, 0}
	ball := ball{pos{float32(winW / 2), float32(winH / 2)}, 20, 400, 400, color{200, 200, 150}}

	keyState := sdl.GetKeyboardState()
	var frameStart time.Time
	var elapsedTime float32
	var controller1Axis int16
	var controller2Axis int16

	playerKey1 := key{up: sdl.SCANCODE_UP, down: sdl.SCANCODE_DOWN}
	playerKey2 := key{up: sdl.SCANCODE_1, down: sdl.SCANCODE_2}

	for {
		frameStart = time.Now()
		// Get Events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return

			}
		}

		for i, controller := range controllerHandlers {
			if controller != nil {
				switch i {
				case 0:
					controller1Axis = controller.Axis(sdl.CONTROLLER_AXIS_LEFTY)
				case 1:
					controller2Axis = controller.Axis(sdl.CONTROLLER_AXIS_LEFTY)
				}
			}
		}

		if state == play {
			// Update State
			player1.update(keyState, playerKey1, controller1Axis, elapsedTime)
			player2.update(keyState, playerKey2, controller2Axis, elapsedTime)
			//player2.aiupdate(&ball, elapsedTime)
			ball.update(&player1, &player2, elapsedTime)
		} else if state == start {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 3 || player2.score == 3 {
					player1.score = 0
					player2.score = 0
				}
				state = play
			}

		}

		// Draw State
		clear(pixels)
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		// Render
		tex.Update(nil, pixels, winW*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())
		if elapsedTime < .0167 {
			sdl.Delay(5 - uint32(elapsedTime/1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
	}
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
