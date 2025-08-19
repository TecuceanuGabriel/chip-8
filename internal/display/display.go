package draw

import (
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
