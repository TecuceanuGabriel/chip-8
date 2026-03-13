package main

import (
	"fmt"
	"os"

	"github.com/TecuceanuGabriel/chip-8/internal/debugger"
	"github.com/TecuceanuGabriel/chip-8/internal/system"

	"github.com/gopxl/pixel/pixelgl"
)

func main() {
	var debugMode bool
	var romPath string

	switch len(os.Args) {
	case 2:
		romPath = os.Args[1]
	case 3:
		if os.Args[1] != "--debug" {
			fmt.Println("usage: ./prog [--debug] rom_path")
			os.Exit(1)
		}
		debugMode = true
		romPath = os.Args[2]
	default:
		fmt.Println("usage: ./prog [--debug] rom_path")
		os.Exit(1)
	}

	sys := system.CreateSystem(romPath, debugMode)
	if debugMode {
		go debugger.Start(sys)
	}
	pixelgl.Run(sys.Run)
}
