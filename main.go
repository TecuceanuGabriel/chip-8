package main

import (
	"fmt"
	"os"

	"github.com/TecuceanuGabriel/chip-8/internal/system"

	"github.com/gopxl/pixel/pixelgl"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./prog rom_path")
		os.Exit(1)
	}

	sys := system.CreateSystem()
	pixelgl.Run(sys.Run)
}
