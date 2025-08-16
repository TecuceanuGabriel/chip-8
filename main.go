package main

import (
	"fmt"
	"os"
)

type System struct {
	memory []byte
	pc     uint16
}

const fontPath = "./FONT"

func load_font() []byte {
	font, err := os.ReadFile(fontPath)
	if err != nil {
		fmt.Printf("Failed to read rom: %v\n", fontPath)
		os.Exit(1)
	}
	return font
}

func fetch(rom []byte) []byte {
	// TODO:
	return rom
}

func decode(instruction []byte) (bool, error) {
	switch instruction {
	//TODO:
	}
	return true, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./prog rom_path")
		os.Exit(1)
	}

	rom_path := os.Args[1]

	rom, err := os.ReadFile(rom_path)
	if err != nil {
		fmt.Printf("Failed to read rom: %v\n", rom_path)
		os.Exit(1)
	}

	fmt.Println(rom[0])

	for {
		instruction := fetch(rom)
		exit, err := decode(instruction)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if exit {
			break
		}
	}
}
