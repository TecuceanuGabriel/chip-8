package main

import (
	"fmt"
	"os"

	"github.com/TecuceanuGabriel/chip-8/internal/system"

	"github.com/faiface/pixel/pixelgl"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./prog rom_path")
		os.Exit(1)
	}

	sys := system.CreateSystem()
	defer sys.Close()

	pixelgl.Run(sys.Run)
}
