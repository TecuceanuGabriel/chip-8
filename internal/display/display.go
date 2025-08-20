package display

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	width  = 64
	height = 32
)

type Display struct {
	pixels [width * height]bool
}

func GetTerminalSize() (width, height int, err error) {
	fd := int(os.Stdout.Fd())
	return term.GetSize(fd)
}

func (d *Display) ClearScreen() {
	fmt.Print("\x1b[2J\x1b[H")
	for i := range height {
		for j := range width {
			d.pixels[j+i*width] = false
		}
	}
}

func (d *Display) DrawSprite(sprite []byte, pos_x, pos_y, n byte) (erasing bool, err error) {
	erasing = false

	tWidth, tHeight, err := GetTerminalSize()
	if err != nil {
		fmt.Println(err)
		return erasing, err
	}

	// TODO: check only on resize
	if tWidth < width || tHeight < height {
		fmt.Println("Terminal window to small")
		os.Exit(1)
	}

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
		var color string
		if d.pixels[idx] == true {
			d.pixels[idx] = false
			color = "255"
			erasing = true
		} else {
			d.pixels[idx] = true
			color = "0"
		}

		// TODO: use one print
		fmt.Printf("\x1b[%v;%vH", pos_y+1, pos_x+1)
		fmt.Printf("\x1b[48;5;%vm", color)
		fmt.Print("█")
		fmt.Print("\x1b[0m") // reset
	}

	return erasing
}
