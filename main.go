package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./prog rom_path")
		os.Exit(1)
	}

	rom_path := os.Args[1]

	rom, err := os.ReadFile(rom_path)
	if err != nil {
		fmt.Println("Failed to read rom: {}", rom_path)
		os.Exit(1)
	}

	fmt.Println(rom[0])
}
