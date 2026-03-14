// Package disasm implements a CHIP-8 disassembler, converting raw opcodes to readable mnemonics.
package disasm

import "fmt"

// Disassemble returns a human-readable string for the two-byte CHIP-8
// instruction. Unknown opcodes are rendered as "UNKNOWN 0xNNNN".
func Disassemble(instr []byte) string {
	b0 := instr[0]
	b1 := instr[1]

	n0 := b0 >> 4                    // first nibble
	x := b0 & 0x0F                   // second nibble (Vx register)
	y := b1 >> 4                     // third nibble (Vy register)
	n := b1 & 0x0F                   // fourth nibble
	kk := b1                         // lower byte
	nnn := uint16(x)<<8 | uint16(b1) // 12-bit address

	switch n0 {
	case 0x0:
		switch b1 {
		case 0xE0:
			return "CLS"
		case 0xEE:
			return "RET"
		default:
			return fmt.Sprintf("SYS  0x%03X", nnn)
		}
	case 0x1:
		return fmt.Sprintf("JP   0x%03X", nnn)
	case 0x2:
		return fmt.Sprintf("CALL 0x%03X", nnn)
	case 0x3:
		return fmt.Sprintf("SE   V%X, 0x%02X", x, kk)
	case 0x4:
		return fmt.Sprintf("SNE  V%X, 0x%02X", x, kk)
	case 0x5:
		return fmt.Sprintf("SE   V%X, V%X", x, y)
	case 0x6:
		return fmt.Sprintf("LD   V%X, 0x%02X", x, kk)
	case 0x7:
		return fmt.Sprintf("ADD  V%X, 0x%02X", x, kk)
	case 0x8:
		switch n {
		case 0x0:
			return fmt.Sprintf("LD   V%X, V%X", x, y)
		case 0x1:
			return fmt.Sprintf("OR   V%X, V%X", x, y)
		case 0x2:
			return fmt.Sprintf("AND  V%X, V%X", x, y)
		case 0x3:
			return fmt.Sprintf("XOR  V%X, V%X", x, y)
		case 0x4:
			return fmt.Sprintf("ADD  V%X, V%X", x, y)
		case 0x5:
			return fmt.Sprintf("SUB  V%X, V%X", x, y)
		case 0x6:
			return fmt.Sprintf("SHR  V%X", x)
		case 0x7:
			return fmt.Sprintf("SUBN V%X, V%X", x, y)
		case 0xE:
			return fmt.Sprintf("SHL  V%X", x)
		}
	case 0x9:
		return fmt.Sprintf("SNE  V%X, V%X", x, y)
	case 0xA:
		return fmt.Sprintf("LD   I, 0x%03X", nnn)
	case 0xB:
		return fmt.Sprintf("JP   V0, 0x%03X", nnn)
	case 0xC:
		return fmt.Sprintf("RND  V%X, 0x%02X", x, kk)
	case 0xD:
		return fmt.Sprintf("DRW  V%X, V%X, %d", x, y, n)
	case 0xE:
		switch b1 {
		case 0x9E:
			return fmt.Sprintf("SKP  V%X", x)
		case 0xA1:
			return fmt.Sprintf("SKNP V%X", x)
		}
	case 0xF:
		switch b1 {
		case 0x07:
			return fmt.Sprintf("LD   V%X, DT", x)
		case 0x0A:
			return fmt.Sprintf("LD   V%X, K", x)
		case 0x15:
			return fmt.Sprintf("LD   DT, V%X", x)
		case 0x18:
			return fmt.Sprintf("LD   ST, V%X", x)
		case 0x1E:
			return fmt.Sprintf("ADD  I, V%X", x)
		case 0x29:
			return fmt.Sprintf("LD   F, V%X", x)
		case 0x33:
			return fmt.Sprintf("LD   B, V%X", x)
		case 0x55:
			return fmt.Sprintf("LD   [I], V%X", x)
		case 0x65:
			return fmt.Sprintf("LD   V%X, [I]", x)
		}
	}

	return fmt.Sprintf("UNKNOWN 0x%02X%02X", b0, b1)
}
