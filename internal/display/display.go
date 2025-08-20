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

func ClearScreen() {
	fmt.Print("\x1b[2J\x1b[H")
}

func (d *Display) drawCell(fill bool, pos_x, pos_y int) {
	idx := pos_x + pos_y*width
	if d.pixels[idx] != fill {
		d.pixels[idx] = !d.pixels[idx]
		color := "0"
		if d.pixels[idx] {
			color = "255"
		}
		fmt.Printf("\x1b[%v;%vH", pos_y, pos_x)
		fmt.Printf("\x1b[48;5;%vm", color)
		fmt.Print(" ")
		fmt.Print("\x1b[0m") // reset
	}
}

func (d *Display) DrawSprite(sprite []byte, pos_x, pos_y, n int) (err error) {
	tWidth, tHeight, err := GetTerminalSize()
	if err != nil {
		fmt.Println(err)
		return err
	}

	if tWidth < width || tHeight < height {
		fmt.Println("Terminal window to small")
		os.Exit(1)
	}

	pos_x = pos_x % width
	pos_y = pos_y % height

	for i := range n {
		line := sprite[i]
		for j := range 8 {
			fill := (line>>j)&1 == 1
			d.drawCell(fill, pos_x+j, pos_y+i)
		}
	}

	return nil
}
