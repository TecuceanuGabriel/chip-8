package display

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	width  = 64
	height = 32
	scale  = 20
)

// Display wraps a pixel window and manages the 64×32 monochrome pixel buffer.
type Display struct {
	pixels [width * height]bool
	win    *pixelgl.Window
}

// NewDisplay opens the application window and returns a ready-to-use Display.
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

// ClearScreen blanks the window and resets all pixels to off.
func (d *Display) ClearScreen() {
	d.win.Clear(colornames.Black)
	for idx := range d.pixels {
		d.pixels[idx] = false
	}
}

// DrawSprite XOR-blits an n-row sprite at (posX, posY), wrapping at screen
// edges. Returns true if any pixel was toggled from on to off (collision).
func (d *Display) DrawSprite(sprite []byte, posX, posY, n byte) (collision bool, err error) {
	collision = false

	// wrap coordinates
	posX = posX % width
	posY = posY % height

	for i := range n {
		if posY+i >= height {
			break
		}

		line := sprite[i]
		for j := range byte(8) {
			if posX+j >= width {
				break
			}

			fill := (line>>(7-j))&1 == 1

			if d.setCell(fill, posX+j, posY+i) {
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

// GetWindow returns the underlying pixel window, used by the system for input
// polling and window-close detection.
func (d *Display) GetWindow() *pixelgl.Window {
	return d.win
}
