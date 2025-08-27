package display

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	width  = 64
	height = 32
	scale  = 20
)

type Display struct {
	pixels [width * height]bool
	win    *pixelgl.Window
}

func NewDisplay() (*Display, error) {
	cfg := pixelgl.WindowConfig{
		Title:  "CHIPY",
		Bounds: pixel.R(0, 0, width*scale, height*scale),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return nil, err
	}

	return &Display{
		win: win,
	}, nil
}

func (d *Display) ClearScreen() {
	d.win.Clear(colornames.Black)
	for y := range height {
		for x := range width {
			d.pixels[getIdx(byte(x), byte(y))] = false
		}
	}
}

func (d *Display) DrawSprite(sprite []byte, pos_x, pos_y, n byte) (collision bool, err error) {
	collision = false

	// wrap coordinates
	pos_x = pos_x % width
	pos_y = pos_y % height

	for i := range n {
		if pos_y+i >= height {
			break
		}

		line := sprite[i]
		for j := range byte(8) {
			if pos_x+j >= width {
				break
			}

			fill := (line>>(7-j))&1 == 1

			if d.setCell(fill, pos_x+j, pos_y+i) {
				collision = true
			}
		}
	}

	d.render()
	return collision, nil
}

func (d *Display) setCell(fill bool, x, y byte) (collision bool) {
	collision = false
	idx := getIdx(x, y)

	if fill {
		if d.pixels[idx] {
			collision = true
		}
		d.pixels[idx] = !d.pixels[idx]
	}

	return collision
}

func (d *Display) render() {
	d.win.Clear(colornames.Black)

	imd := imdraw.New(nil)
	imd.Color = colornames.White

	for y := range height {
		for x := range width {
			idx := getIdx(byte(x), byte(y))
			if d.pixels[idx] {
				min := pixel.V(float64(x)*scale, float64(height-int(y))*scale)
				max := pixel.V(float64(x+1)*scale, float64(height-int(y+1))*scale)
				imd.Push(min, max)
				imd.Rectangle(0)
			}
		}
	}

	imd.Draw(d.win)
}

func getIdx(x, y byte) uint16 {
	return uint16(x) + uint16(y)*width
}

func (d *Display) GetWindow() *pixelgl.Window {
	return d.win
}
