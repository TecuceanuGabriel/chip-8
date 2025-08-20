package display

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const (
	width  = 64
	height = 32
)

type Display struct {
	pixels [width * height]bool
}

func GetTerminalSize() (width, height int, err error) {
	cmd := exec.Command("stty", "size")
	out_str, err := cmd.Output()

	if err != nil {
		return 0, 0, err
	}

	fields := strings.Fields(strings.Trim(string(out_str), " "))

	width, err = strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, err
	}

	height, err = strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
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
	width, height, err := GetTerminalSize()
	if err != nil {
		return err
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

	d.drawCell(true, 10, 10)

	return nil
}
