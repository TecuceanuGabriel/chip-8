// config.go loads chip8.toml and exposes the Quirks and keymap configuration.
package system

import (
	"fmt"
	"unicode"

	"github.com/BurntSushi/toml"
)

const configPath = "chip8.toml"

// Quirks controls the three instructions whose behaviour differs between
// original CHIP-8 and CHIP-48/SUPER-CHIP.
type Quirks struct {
	// ShiftUsesVY: 8XY6/8XYE shift VY and store in VX (true=original CHIP-8).
	// When false (CHIP-48 default) VX is shifted in place and VY is ignored.
	ShiftUsesVY bool `toml:"shift_uses_vy"`

	// JumpUsesVX: BNNN uses VX instead of V0 as the base register (true=CHIP-48).
	// When false (original default) JP V0, NNN is used.
	JumpUsesVX bool `toml:"jump_uses_vx"`

	// LoadStoreIncI: FX55/FX65 increment I after each register access (true=original CHIP-8).
	// When false (CHIP-48 default) I is left unchanged.
	LoadStoreIncI bool `toml:"load_store_inc_i"`
}

type keymapConfig struct {
	Layout []string `toml:"layout"`
}

type fileConfig struct {
	Keymap keymapConfig `toml:"keymap"`
	Quirks Quirks       `toml:"quirks"`
}

// loadConfig reads chip8.toml from the working directory.
// Missing file or missing sections fall back to sane defaults silently.
func loadConfig() (map[byte]byte, Quirks) {
	var cfg fileConfig
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		// File absent or unreadable — use defaults without noise.
		return defaultKeymap(), Quirks{}
	}

	km := defaultKeymap()
	if len(cfg.Keymap.Layout) == 4 {
		if parsed, ok := parseKeymapLayout(cfg.Keymap.Layout); ok {
			km = parsed
		} else {
			fmt.Println("warning: chip8.toml keymap.layout must be 4 rows of 4 characters; using default")
		}
	}

	return km, cfg.Quirks
}

func defaultKeymap() map[byte]byte {
	km, _ := parseKeymapLayout([]string{"1234", "QWER", "ASDF", "ZXCV"})
	return km
}

func parseKeymapLayout(layout []string) (map[byte]byte, bool) {
	chip8Keys := []byte{
		0x1, 0x2, 0x3, 0xC,
		0x4, 0x5, 0x6, 0xD,
		0x7, 0x8, 0x9, 0xE,
		0xA, 0x0, 0xB, 0xF,
	}

	km := make(map[byte]byte, 16)
	i := 0
	for _, row := range layout {
		runes := []rune(row)
		if len(runes) != 4 {
			return nil, false
		}
		for _, ch := range runes {
			km[byte(unicode.ToUpper(ch))] = chip8Keys[i]
			i++
		}
	}
	if i != 16 {
		return nil, false
	}
	return km, true
}
