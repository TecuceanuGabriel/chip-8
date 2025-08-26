package display

import (
	"image/color"

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
	sprite *pixel.Sprite
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
	for i := range height {
		for j := range width {
			d.pixels[j+i*width] = false
		}
	}
}

func (d *Display) DrawSprite(sprite []byte, pos_x, pos_y, n byte) (erasing bool, err error) {
	erasing = false

	// wrap coordinates
	pos_x = pos_x % width
	pos_y = pos_y % height

	for i := range n {
		line := sprite[i]
		for j := range byte(8) {
			fill := (line>>(7-j))&1 == 1
			e := d.drawCell(fill, pos_x+j, pos_y+i)
			if e {
				erasing = e
			}
		}
	}

	return erasing, nil
}

func (d *Display) drawCell(fill bool, pos_x, pos_y byte) (erasing bool) {
	erasing = false
	idx := pos_x + pos_y*width

	if fill == true {
		if d.pixels[idx] == true {
			d.pixels[idx] = false
			erasing = true
		} else {
			d.pixels[idx] = true
		}

		min := pixel.V(float64(pos_x)*scale, float64(height-int(pos_y))*scale)
		max := pixel.V(float64(pos_x+1)*scale, float64(height-int(pos_y+1))*scale)

		imd := imdraw.New(nil)
		imd.Color = color.White
		imd.Push(min, max)
		imd.Rectangle(0)
		imd.Draw(d.win)
	}

	return erasing
}

func (d *Display) GetWindow() *pixelgl.Window {
	return d.win
}
