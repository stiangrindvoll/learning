package main

import (
	"bufio"
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type field struct {
	squares  squares
	selected int
}

type square struct {
	R        sdl.Rect
	Selected bool
}

// ID of the square
type ID int32

type squares []square

func createSquares(startX, startY, width, height, spacing int32) (sq squares) {
	var x, y int32

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			x = (int32(j) * (spacing + width)) + startX
			y = (int32(i) * (spacing + height)) + startY
			sq = append(sq, square{R: sdl.Rect{X: x, Y: y, W: width, H: height}})

		}
	}
	return
}

func (f *field) render(r *sdl.Renderer, mp sdl.Point) {
	for i, s := range f.squares {
		if f.selected >= 0 && f.selected < 9 && f.selected == i {
			r.SetDrawColor(100, 255, 255, 255)
		} else if mp.InRect(&s.R) {
			r.SetDrawColor(255, 0, 255, 255)
		} else {
			r.SetDrawColor(100, 0, 255, 255)
		}

		r.FillRect(&s.R)
	}
}

func (f *field) setSelected(mp sdl.Point, communicataion *bufio.ReadWriter) {
	for i, s := range f.squares {
		if mp.InRect(&s.R) {
			f.selected = i
			communicataion.WriteString(fmt.Sprintf("selected: %d\n", i))
			fmt.Println("Selected Square:", i)
		}
	}
}
