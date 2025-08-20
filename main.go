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
