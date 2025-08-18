package main

import (
	"fmt"
	"os"
)

type System struct {
	memory    []byte
	pc        uint16
	registers []byte
}

const (
	memorySize          = 4096
	firstInstructionAdd = 0x200
	fontStartAdd        = 0x50
)

var font = []byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func createSystem() *System {
	system := System{
		memory:    make([]byte, memorySize),
		pc:        firstInstructionAdd,
		registers: make([]byte, 16),
	}

	copy(system.memory[fontStartAdd:], font)

	rom_path := os.Args[1]

	rom, err := os.ReadFile(rom_path)
	if err != nil {
		fmt.Printf("Failed to read rom: %v\n", rom_path)
		os.Exit(1)
	}

	copy(system.memory[firstInstructionAdd:], rom)

	return &system
}

func fetch(system *System) []byte {
	value := system.memory[system.pc : system.pc+2]
	system.pc += 2
	return value
}

func decode(instruction []byte) (bool, error) {
	fmt.Printf("Instruction: %x\n", instruction)

	firstByte := instruction[0]
	secondByte := instruction[1]

	firstNibble := firstByte >> 4
	secondNibble := firstByte & 0x0F
	thirdNibble := secondByte >> 4
	fourthNibble := secondByte & 0x0F

	last3Nibbles := (uint16(secondNibble) << 8) | (uint16(thirdNibble) << 4) | uint16(fourthNibble)

	fmt.Printf("Bytes: %x%x\n", firstByte, secondByte)
	fmt.Printf("nibbles: %x%x%x%x\n", firstNibble, secondNibble, thirdNibble, fourthNibble)
	fmt.Printf("last 3: %x\n", last3Nibbles)

	switch firstNibble {

	}

	return true, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./prog rom_path")
		os.Exit(1)
	}

	system := createSystem()

	for {
		instruction := fetch(system)
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
