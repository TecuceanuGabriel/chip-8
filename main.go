package main

import (
	"fmt"
	"os"

	"github.com/TecuceanuGabriel/chip-8/internal/display"
	"github.com/TecuceanuGabriel/chip-8/internal/stack"
)

type System struct {
	memory     []byte
	pc         uint16
	call_stack stack.Stack[[2]byte]
	registers  []byte
	iReg       uint16
	display    display.Display
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

func createSystem() (system *System) {
	system = &System{
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

	return system
}

func (system *System) fetch() (instruction []byte) {
	instruction = system.memory[system.pc : system.pc+2]
	system.pc += 2
	return instruction
}

func (system *System) decode(instruction []byte) (bool, error) {
	// fmt.Printf("Instruction: %x\n", instruction)

	firstByte := instruction[0]
	secondByte := instruction[1]

	firstNibble := firstByte >> 4
	secondNibble := firstByte & 0x0F
	thirdNibble := secondByte >> 4
	fourthNibble := secondByte & 0x0F

	last2Nibbles := (thirdNibble << 4) | fourthNibble
	last3Nibbles := (uint16(secondNibble) << 8) | uint16(last2Nibbles)

	// fmt.Printf("Bytes: %x%x\n", firstByte, secondByte)
	// fmt.Printf("nibbles: %x%x%x%x\n", firstNibble, secondNibble, thirdNibble, fourthNibble)
	// fmt.Printf("last 2: %x\n", last2Nibbles)
	// fmt.Printf("last 3: %x\n", last3Nibbles)

	switch firstNibble {
	case 0:
		{
			if fourthNibble == 0 {
				system.display.ClearScreen()
			} else {
				// TODO: return from subrutine
			}
		}
	case 1:
		{
			// TODO: jump to nnn
			system.pc = last3Nibbles
		}
	case 6:
		{
			// TODO: load nn to reg
			system.registers[secondNibble] = last2Nibbles
		}
	case 7:
		{
			// TODO: add to reg nn
			system.registers[secondNibble] += last2Nibbles
		}
	case 0xA:
		{
			// TODO: set index reg I
			system.iReg = last3Nibbles
		}
	case 0xD:
		{
			// TODO: draw
			sprite := system.memory[system.iReg : system.iReg+uint16(fourthNibble)]
			system.registers[0xF] = 0
			pos_x := system.registers[secondNibble]
			pos_y := system.registers[thirdNibble]
			erasing, err := system.display.DrawSprite(sprite, pos_x, pos_y, fourthNibble)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if erasing {
				system.registers[0xF] = 1
			}
		}

	}

	return false, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./prog rom_path")
		os.Exit(1)
	}

	system := createSystem()

	for {
		instruction := system.fetch()
		exit, err := system.decode(instruction)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if exit {
			break
		}
	}
}
