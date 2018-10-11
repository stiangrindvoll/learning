package main

import "github.com/veandco/go-sdl2/sdl"

type field struct {
	squares  squares
	selected *square
}

type square struct {
	ID       ID
	R        sdl.Rect
	Selected bool
}

// ID of the square
type ID int32

type squares []square

func createSquares(startX, startY, width, height, spacing int32) (sq squares) {
	var x, y int32
	var id ID

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			id++
			x = (int32(i) * (spacing + width)) + startX
			y = (int32(j) * (spacing + height)) + startY
			sq = append(sq, square{R: sdl.Rect{X: x, Y: y, W: width, H: height}, ID: id})

		}
	}
	return
}

func (sq squares) render(r *sdl.Renderer, mp sdl.Point) {
	for _, s := range sq {
		if mp.InRect(&s.R) {
			r.SetDrawColor(255, 0, 255, 255)
		} else if s.Selected {
			r.SetDrawColor(100, 255, 255, 255)

		} else {
			r.SetDrawColor(100, 0, 255, 255)
		}

		r.FillRect(&s.R)
	}
}

func (sq squares) setSelected(mp sdl.Point) {
	for i, s := range sq {
		if mp.InRect(&s.R) {
			sq[i].Selected = true
		} else {
			sq[i].Selected = false
		}
	}
}
