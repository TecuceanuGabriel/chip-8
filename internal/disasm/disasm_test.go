package disasm

import "testing"

func TestDisassemble(t *testing.T) {
	tests := []struct {
		instr [2]byte
		want  string
	}{
		// 0x0 — system
		{[2]byte{0x00, 0xE0}, "CLS"},
		{[2]byte{0x00, 0xEE}, "RET"},
		// 0x1–0x7 — control flow and immediate ops
		{[2]byte{0x12, 0xAB}, "JP   0x2AB"},
		{[2]byte{0x22, 0xAB}, "CALL 0x2AB"},
		{[2]byte{0x31, 0xFF}, "SE   V1, 0xFF"},
		{[2]byte{0x41, 0xFF}, "SNE  V1, 0xFF"},
		{[2]byte{0x51, 0x20}, "SE   V1, V2"},
		{[2]byte{0x61, 0x42}, "LD   V1, 0x42"},
		{[2]byte{0x71, 0x01}, "ADD  V1, 0x01"},
		// 0x8 — arithmetic
		{[2]byte{0x81, 0x20}, "LD   V1, V2"},
		{[2]byte{0x81, 0x21}, "OR   V1, V2"},
		{[2]byte{0x81, 0x22}, "AND  V1, V2"},
		{[2]byte{0x81, 0x23}, "XOR  V1, V2"},
		{[2]byte{0x81, 0x24}, "ADD  V1, V2"},
		{[2]byte{0x81, 0x25}, "SUB  V1, V2"},
		{[2]byte{0x81, 0x06}, "SHR  V1"},
		{[2]byte{0x81, 0x27}, "SUBN V1, V2"},
		{[2]byte{0x81, 0x2E}, "SHL  V1"},
		// 0x9–0xD
		{[2]byte{0x91, 0x20}, "SNE  V1, V2"},
		{[2]byte{0xA2, 0xAB}, "LD   I, 0x2AB"},
		{[2]byte{0xB2, 0xAB}, "JP   V0, 0x2AB"},
		{[2]byte{0xC1, 0xFF}, "RND  V1, 0xFF"},
		{[2]byte{0xD1, 0x23}, "DRW  V1, V2, 3"},
		// 0xE — key
		{[2]byte{0xE1, 0x9E}, "SKP  V1"},
		{[2]byte{0xE1, 0xA1}, "SKNP V1"},
		// 0xF — misc
		{[2]byte{0xF1, 0x07}, "LD   V1, DT"},
		{[2]byte{0xF1, 0x0A}, "LD   V1, K"},
		{[2]byte{0xF1, 0x15}, "LD   DT, V1"},
		{[2]byte{0xF1, 0x18}, "LD   ST, V1"},
		{[2]byte{0xF1, 0x1E}, "ADD  I, V1"},
		{[2]byte{0xF1, 0x29}, "LD   F, V1"},
		{[2]byte{0xF1, 0x33}, "LD   B, V1"},
		{[2]byte{0xF1, 0x55}, "LD   [I], V1"},
		{[2]byte{0xF1, 0x65}, "LD   V1, [I]"},
		// unknown
		{[2]byte{0xFF, 0xFF}, "UNKNOWN 0xFFFF"},
	}

	for _, tt := range tests {
		got := Disassemble(tt.instr[:])
		if got != tt.want {
			t.Errorf("Disassemble(%02X %02X) = %q, want %q",
				tt.instr[0], tt.instr[1], got, tt.want)
		}
	}
}
