package debugger

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/TecuceanuGabriel/chip-8/internal/disasm"
	"github.com/TecuceanuGabriel/chip-8/internal/system"
	"github.com/chzyer/readline"
)

const contextLines = 5 // instructions to show around PC on each pause

// Start runs the debugger REPL. It is intended to be called as a goroutine
// before pixelgl.Run so it runs concurrently with the game loop.
func Start(sys *system.System) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "(dbg) ",
		HistoryFile: "/tmp/chip8_debug_history",
	})
	if err != nil {
		fmt.Printf("debugger: failed to init readline: %v\n", err)
		return
	}
	defer rl.Close()

	// The emulator starts paused; show context immediately.
	printContext(sys)

	for {
		line, err := rl.Readline()
		if err == io.EOF || err == readline.ErrInterrupt {
			sys.DebugChan() <- system.CmdQuit{}
			return
		}
		if err != nil {
			fmt.Printf("debugger: readline error: %v\n", err)
			return
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "s", "step":
			count := 1
			if len(parts) > 1 {
				n, err := strconv.Atoi(parts[1])
				if err != nil || n <= 0 {
					fmt.Println("usage: step [N]")
					continue
				}
				count = n
			}
			sys.DebugChan() <- system.CmdStep{Count: count}
			<-sys.EventChan() // wait for the step to complete
			printContext(sys)

		case "c", "continue":
			sys.DebugChan() <- system.CmdContinue{}
			<-sys.EventChan() // block until the next breakpoint or manual break.
			printContext(sys)

		case "q", "quit":
			sys.DebugChan() <- system.CmdQuit{}
			return

		default:
			fmt.Printf("unknown command %q — try: step [N], continue, quit\n", parts[0])
		}
	}
}

// printContext disassembles contextLines instructions starting from the current PC.
func printContext(sys *system.System) {
	pc := sys.PC()
	fmt.Println()
	for i := range contextLines {
		addr := pc + uint16(i*2)
		mem := sys.MemorySlice(addr, 2)
		if len(mem) < 2 {
			break
		}
		marker := "  "
		if addr == pc {
			marker = "→ "
		}
		fmt.Printf("%s0x%03X  %s\n", marker, addr, disasm.Disassemble(mem))
	}
	fmt.Println()
}
