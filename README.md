A chip-8 interpreter written in go.

# Usage

```bash
go run . path/to/rom.ch8
```

|  Key  | Action          |
| :---: | --------------- |
| Space | Pause execution |
|  Esc  | Exit            |

# Features

- complete CHIP-8 instruction set (all original opcodes implemented).
- gui rendering using gopxl/pixel.
- sound support.
- keyboard input mapping (default in ./KEYMAP).

# TODOs

- [ ] add configuration options (execution speed, scaling, etc.).
- [ ] add debug mode.
- [ ] add prebuilt binaries.
