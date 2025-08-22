package system

import (
	"fmt"
	"math/rand/v2"
	"os"
	"syscall"
	"time"
	"unicode"

	"github.com/eiannone/keyboard"

	"github.com/TecuceanuGabriel/chip-8/internal/display"
	"github.com/TecuceanuGabriel/chip-8/internal/stack"
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

const (
	memorySize                 = 4096
	firstInstructionAdd uint16 = 0x200
	fontStartAddr       uint16 = 0x50
	keymapPath                 = "./KEYMAP"
)

type System struct {
	memory     []byte
	pc         uint16
	call_stack stack.Stack[uint16]
	registers  []byte
	iReg       uint16
	display    display.Display

	keymap  map[byte]byte
	keyChan <-chan keyboard.KeyEvent

	soundTimer byte
	delayTimer byte
}

func CreateSystem() (system *System) {
	keyChan, err := keyboard.GetKeys(10)
	if err != nil {
		fmt.Println("Failed to create keyboard channel", err)
		os.Exit(1)
	}

	system = &System{
		memory:    make([]byte, memorySize),
		pc:        firstInstructionAdd,
		registers: make([]byte, 16),
		keymap:    loadKeymap(),
		keyChan:   keyChan,
	}

	copy(system.memory[fontStartAddr:], font)

	rom_path := os.Args[1]

	rom, err := os.ReadFile(rom_path)
	if err != nil {
		fmt.Printf("Failed to load rom: %v\n", rom_path)
		os.Exit(1)
	}

	copy(system.memory[firstInstructionAdd:], rom)

	go system.listenForExit()

	return system
}

func (system *System) listenForExit() {
	for key := range system.keyChan {
		if key.Key == keyboard.KeyCtrlC {
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}
	}
}

func loadKeymap() (keymap map[byte]byte) {
	file, err := os.ReadFile(keymapPath)
	if err != nil {
		fmt.Printf("Failed to load keymap: %v\n", keymapPath)
	}

	original_keys := []byte{
		1, 2, 3, 0xC,
		4, 5, 6, 0xD,
		7, 8, 9, 0xE,
		0xA, 0, 0xB, 0xF,
	}

	keymap = make(map[byte]byte, 16)
	var i int
	for _, b := range file {
		if b == '\n' {
			continue
		}

		key := byte(unicode.ToLower(rune(b)))
		keymap[key] = original_keys[i]
		i++
	}

	return keymap
}

func (system *System) Close() {
	keyboard.Close()
}

func (system *System) HandleTimers() {
	go system.handleDelayTimer()
	go system.handleSoundTimer()
}

func (system *System) handleDelayTimer() {
	if system.delayTimer == 0 {
		return
	}

	executeAt60Hz(func() bool {
		system.delayTimer--
		return system.delayTimer == 0
	})
}

func (system *System) handleSoundTimer() {
	if system.soundTimer == 0 {
		return
	}

	executeAt60Hz(func() bool {
		system.soundTimer--

		if system.soundTimer != 0 {
			fmt.Print("\a") // beep
			return false
		}

		return true
	})
}

func executeAt60Hz(action func() bool) {
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	for range ticker.C {
		stop := action()
		if stop {
			break
		}
	}
}

func (system *System) Fetch() (instruction []byte) {
	instruction = system.memory[system.pc : system.pc+2]
	system.pc += 2
	return instruction
}

func (system *System) Decode(instruction []byte) error {
	firstByte := instruction[0]
	secondByte := instruction[1]

	firstNibble := firstByte >> 4
	secondNibble := firstByte & 0x0F
	thirdNibble := secondByte >> 4
	fourthNibble := secondByte & 0x0F

	last2Nibbles := (thirdNibble << 4) | fourthNibble
	last3Nibbles := (uint16(secondNibble) << 8) | uint16(last2Nibbles)

	switch firstNibble {
	case 0:
		switch secondByte {
		case 0xE0: // CLS
			system.display.ClearScreen()
		case 0xEE: // RET
			system.ret()
		}
	case 1: // JP addr
		system.pc = last3Nibbles
	case 2: // CALL addr
		system.call(last3Nibbles)
	case 3: // SE Vx, byte
		system.skip_equal_im(secondNibble, secondByte)
	case 4: // SNE Vx, byte
		system.skip_not_equal_im(secondNibble, secondByte)
	case 5: // SE Vx, Vy
		system.skip_equal(secondNibble, thirdNibble)
	case 6: // LD Vx, byte
		system.registers[secondNibble] = last2Nibbles
	case 7: // ADD Vx, byte
		system.registers[secondNibble] += last2Nibbles
	case 8:
		system.decodeArithmetic(fourthNibble, secondNibble, thirdNibble)
	case 9: // SNE Vx, Vy
		system.skip_not_equal(secondNibble, thirdNibble)
	case 0xA: // LD I, addr
		system.iReg = last3Nibbles
	case 0xB: // JP V0, addr TODO: comp problem: make configurable?
		system.pc = uint16(system.registers[0]) + last3Nibbles
	case 0xC: // RND Vx, byte
		system.registers[secondNibble] = byte(rand.UintN(256)) & secondByte
	case 0xD: // DRW Vx, Vy, nibble
		system.drw(secondNibble, thirdNibble, fourthNibble)
	case 0xE:
		{
			switch secondByte {
			case 0x9E: // SKP Vx
				system.skip_pressed(secondNibble)
			case 0xA1: // SKNP Vx
				system.skip_not_pressed(secondNibble)
			}
		}
	case 0xF:
		system.decodeF(secondByte, secondNibble)
	default:
		{
			fmt.Printf("Unknown instruction: %x\n", instruction)
			os.Exit(1)
		}
	}

	return nil
}

func (system *System) decodeArithmetic(instType, x_addr, y_addr byte) {
	switch instType {
	case 0: // LD Vx, Vy
		system.registers[x_addr] = system.registers[y_addr]
	case 1: // OR Vx, Vy
		system.registers[x_addr] |= system.registers[y_addr]
	case 2: // AND Vx, Vy
		system.registers[x_addr] &= system.registers[y_addr]
	case 3: // XOR Vx, Vy
		system.registers[x_addr] ^= system.registers[y_addr]
	case 4: // ADD Vx, Vy
		system.add(x_addr, y_addr)
	case 5: // SUB Vx, Vy
		system.sub(x_addr, y_addr)
	case 6: // SHR Vx {, Vy}
		system.shr(x_addr)
	case 7: // SUBN Vx, Vy
		system.subn(x_addr, y_addr)
	case 0xE: // AND Vx, Vy
		system.shl(x_addr)
	default:
		{
			fmt.Printf("Unknown arithmetic instruction: %x\n", instType)
			os.Exit(1)
		}
	}
}

func (system *System) decodeF(instType, x_addr byte) {
	switch instType {
	case 0x07: // LD Vx, DT
		system.registers[x_addr] = system.delayTimer
	case 0x0A: // LD Vx, K
		system.get_key(x_addr)
	case 0x15: // LD DT, Vx
		system.delayTimer = system.registers[x_addr]
	case 0x18: // LD ST, Vx
		system.soundTimer = system.registers[x_addr]
	case 0x1E: // ADD I, Vx
		system.addToIReg(x_addr)
	case 0x29: // LD F, Vx
		system.setFontLoc(x_addr)
	default:
		{
			fmt.Printf("Unknown arithmetic instruction: %x\n", instType)
			os.Exit(1)
		}
	}
}

func (system *System) ret() {
	old_pc, err := system.call_stack.Pop()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	system.pc = old_pc
}

func (system *System) call(addr uint16) {
	system.call_stack.Push(system.pc)
	system.pc = addr
}

func (system *System) skip_equal(x_addr, y_addr byte) {
	if system.registers[x_addr] == system.registers[y_addr] {
		system.pc += 2
	}
}

func (system *System) skip_not_equal(x_addr, y_addr byte) {
	if system.registers[x_addr] != system.registers[y_addr] {
		system.pc += 2
	}
}

func (system *System) skip_equal_im(x_addr, val byte) {
	if system.registers[x_addr] == val {
		system.pc += 2
	}
}

func (system *System) skip_not_equal_im(x_addr, val byte) {
	if system.registers[x_addr] != val {
		system.pc += 2
	}
}

func (system *System) add(x_addr, y_addr byte) {
	result := uint16(system.registers[x_addr]) + uint16(system.registers[y_addr])
	if result > 255 {
		system.registers[0xF] = 1
	} else {
		system.registers[0xF] = 0
	}
	system.registers[x_addr] = byte(result)
}

func (system *System) sub(x_addr, y_addr byte) {
	x := system.registers[x_addr]
	y := system.registers[y_addr]
	if x > y {
		system.registers[0xF] = 1
	} else {
		system.registers[0xF] = 0
	}
	system.registers[x_addr] = x - y
}

func (system *System) subn(x_addr, y_addr byte) {
	x := system.registers[x_addr]
	y := system.registers[y_addr]
	if y > x {
		system.registers[0xF] = 1
	} else {
		system.registers[0xF] = 0
	}
	system.registers[x_addr] = y - x
}

// TODO: make shifts configurable (use Vy or not)
func (system *System) shr(x_addr byte) {
	x := system.registers[x_addr]
	if x&1 == 1 {
		system.registers[0xF] = 1
	} else {
		system.registers[0xF] = 0
	}
	system.registers[x_addr] >>= 1
}

func (system *System) shl(x_addr byte) {
	x := system.registers[x_addr]
	if (x>>7)&1 == 1 {
		system.registers[0xF] = 1
	} else {
		system.registers[0xF] = 0
	}
	system.registers[x_addr] <<= 1
}

func (system *System) drw(x_addr, y_addr, n byte) {
	sprite := system.memory[system.iReg : system.iReg+uint16(n)]

	pos_x := system.registers[x_addr]
	pos_y := system.registers[y_addr]

	system.registers[0xF] = 0
	erasing, err := system.display.DrawSprite(sprite, pos_x, pos_y, n)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if erasing {
		system.registers[0xF] = 1
	}
}

func (system *System) skip_pressed(x_addr byte) {
	key, hasInput := system.GetPressedKey()
	if hasInput && key == system.registers[x_addr] {
		system.pc += 2
	}
}

func (system *System) skip_not_pressed(x_addr byte) {
	key, hasInput := system.GetPressedKey()
	if !hasInput || (hasInput && (key != system.registers[x_addr])) {
		system.pc += 2
	}
}

func (system *System) get_key(x_addr byte) {
	key, hasInput := system.GetPressedKey()
	if hasInput {
		system.registers[x_addr] = key
	} else {
		system.pc -= 2 // block
	}
}

func (system *System) GetPressedKey() (byte, bool) {
	select {
	case key := <-system.keyChan:
		return system.keymap[byte(key.Rune)], true
	default:
		return 0, false
	}
}

func (system *System) addToIReg(x_addr byte) {
	result := system.iReg + uint16(system.registers[x_addr])
	if result > 0x0FFF { // only left-most 12 bits are used
		system.registers[0xF] = 1
	} else {
		system.registers[0xF] = 0
	}
	system.iReg = result
}

func (system *System) setFontLoc(x_addr byte) {
	char := system.registers[x_addr] & 0x0F
	pos := fontStartAddr + uint16(char*5)
	system.iReg = pos
}
