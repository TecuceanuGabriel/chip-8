package main

import (
	"fmt"
	"os"

	"github.com/TecuceanuGabriel/chip-8/internal/system"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./prog rom_path")
		os.Exit(1)
	}

	system := system.CreateSystem()

	initTerminal()
	defer restoreTerminal()

	for {
		instruction := system.Fetch()
		exit, err := system.Decode(instruction)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if exit {
			break
		}
	}
}

func initTerminal() {
	fmt.Print("\x1b[?25l") // hide cursor
}

func restoreTerminal() {
	// TODO: print only once
	fmt.Print("\x1b[0m")   // reset all attributes (colors, styles)
	fmt.Print("\x1b[?25h") // show cursor
	fmt.Print("\x1b[2J")   // clear screen
	fmt.Print("\x1b[H")    // move cursor to home position
}
