// Package debugger provides an interactive REPL for the CHIP-8 debugger.
package debugger

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/TecuceanuGabriel/chip-8/internal/disasm"
	"github.com/TecuceanuGabriel/chip-8/internal/system"
	"github.com/charmbracelet/lipgloss"
	"github.com/chzyer/readline"
	"golang.org/x/term"
)

const contextLines = 5 // instructions shown around PC on each pause

// Start runs the debugger REPL. It is intended to be called as a goroutine
// before pixelgl.Run so it runs concurrently with the game loop.
func Start(sys *system.System) {
	keyItems := make([]readline.PrefixCompleterInterface, 16)
	for i := range 16 {
		keyItems[i] = readline.PcItem(fmt.Sprintf("%X", i))
	}

	completer := readline.NewPrefixCompleter(
		readline.PcItem("step"),
		readline.PcItem("continue"),
		readline.PcItem("audio"),
		readline.PcItem("reset"),
		readline.PcItem("quit"),
		readline.PcItem("b"),
		readline.PcItem("rb"),
		readline.PcItem("lb"),
		readline.PcItem("regs"),
		readline.PcItem("mem"),
		readline.PcItem("timers"),
		readline.PcItem("keys"),
		readline.PcItem("dis"),
		readline.PcItem("press", keyItems...),
		readline.PcItem("release", keyItems...),
	)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "(dbg) ",
		HistoryFile:  "/tmp/chip8_debug_history",
		AutoComplete: completer,
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
		// --- execution control ---

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
			<-sys.EventChan()
			printContext(sys)

		case "c", "continue":
			sys.DebugChan() <- system.CmdContinue{}
			go func() {
				<-sys.EventChan()
				rl.Clean()
				printContext(sys)
				rl.Refresh()
			}()

		case "audio":
			sys.DebugChan() <- system.CmdToggleAudio{}
			if sys.AudioMuted() {
				fmt.Println("audio off")
			} else {
				fmt.Println("audio on")
			}

		case "reset":
			sys.DebugChan() <- system.CmdReset{}
			<-sys.EventChan()
			printContext(sys)

		case "q", "quit":
			sys.DebugChan() <- system.CmdQuit{}
			return

		// --- breakpoints ---

		case "b":
			addr, ok := parseAddr(parts, 1)
			if !ok {
				fmt.Println("usage: b 0xADDR")
				continue
			}
			sys.DebugChan() <- system.CmdSetBreakpoint{Addr: addr}
			fmt.Printf("breakpoint set at 0x%03X\n", addr)

		case "rb":
			addr, ok := parseAddr(parts, 1)
			if !ok {
				fmt.Println("usage: rb 0xADDR")
				continue
			}
			sys.DebugChan() <- system.CmdRemoveBreakpoint{Addr: addr}
			fmt.Printf("breakpoint removed at 0x%03X\n", addr)

		case "lb":
			addrs := sys.Breakpoints()
			if len(addrs) == 0 {
				fmt.Println("no breakpoints set")
			} else {
				sort.Slice(addrs, func(i, j int) bool { return addrs[i] < addrs[j] })
				for _, a := range addrs {
					fmt.Printf("  0x%03X\n", a)
				}
			}

		// --- key injection ---

		case "press":
			key, ok := parseKey(parts, 1)
			if !ok {
				fmt.Println("usage: press <0-9|A-F>")
				continue
			}
			sys.DebugChan() <- system.CmdSetKey{Key: key}
			fmt.Printf("key %X pressed\n", key)

		case "release":
			key, ok := parseKey(parts, 1)
			if !ok {
				fmt.Println("usage: release <0-9|A-F>")
				continue
			}
			sys.DebugChan() <- system.CmdReleaseKey{Key: key}
			fmt.Printf("key %X released\n", key)

		// --- state views ---

		case "r", "regs":
			printRegs(sys)

		case "m", "mem":
			addr, ok := parseAddr(parts, 1)
			if !ok {
				fmt.Println("usage: m 0xADDR [N]")
				continue
			}
			n := 16
			if len(parts) > 2 {
				v, err := strconv.Atoi(parts[2])
				if err != nil || v <= 0 {
					fmt.Println("usage: m 0xADDR [N]")
					continue
				}
				n = v
			}
			printMem(sys, addr, n)

		case "t", "timers":
			fmt.Printf("delay: %d   sound: %d\n", sys.DelayTimer(), sys.SoundTimer())

		case "k", "keys":
			printKeys(sys)

		case "d", "dis":
			n := contextLines
			if len(parts) > 1 {
				v, err := strconv.Atoi(parts[1])
				if err != nil || v <= 0 {
					fmt.Println("usage: dis [N]")
					continue
				}
				n = v
			}
			printDis(sys, sys.PC(), n)

		default:
			fmt.Printf("unknown command %q\n", parts[0])
			fmt.Println("commands: step [N], continue, quit, b/rb/lb, regs, mem, timers, keys, dis [N], press/release <key>")
		}
	}
}

// --- helpers ---

func parseKey(parts []string, idx int) (byte, bool) {
	if len(parts) <= idx || len(parts[idx]) != 1 {
		return 0, false
	}
	v, err := strconv.ParseUint(strings.ToUpper(parts[idx]), 16, 8)
	if err != nil || v > 0xF {
		return 0, false
	}
	return byte(v), true
}

func parseAddr(parts []string, idx int) (uint16, bool) {
	if len(parts) <= idx {
		return 0, false
	}
	s := strings.TrimPrefix(parts[idx], "0x")
	s = strings.TrimPrefix(s, "0X")
	v, err := strconv.ParseUint(s, 16, 16)
	if err != nil {
		return 0, false
	}
	return uint16(v), true
}

var (
	colStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	addrColor   = lipgloss.NewStyle().Foreground(lipgloss.Color("6")) // cyan
	mnemoColor  = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // yellow
	arrowStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	regValColor = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // green
	keyOnStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	keyOffStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func printContext(sys *system.System) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width < 80 {
		printDis(sys, sys.PC(), contextLines)
		printRegs(sys)
		printKeys(sys)
		return
	}

	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top,
		colStyle.Render(disasmBlock(sys)),
		colStyle.Render(regsBlock(sys)),
		colStyle.Render(keysBlock(sys)),
	))
}

func disasmBlock(sys *system.System) string {
	pc := sys.PC()
	var sb strings.Builder
	for i := range contextLines {
		addr := pc + uint16(i*2)
		mem := sys.MemorySlice(addr, 2)
		if len(mem) < 2 {
			break
		}
		if i > 0 {
			sb.WriteByte('\n')
		}
		if addr == pc {
			sb.WriteString(arrowStyle.Render("→ "))
		} else {
			sb.WriteString("  ")
		}
		parts := strings.SplitN(disasm.Disassemble(mem), " ", 2)
		sb.WriteString(addrColor.Render(fmt.Sprintf("0x%03X", addr)) + "  ")
		sb.WriteString(mnemoColor.Render(parts[0]))
		if len(parts) > 1 {
			sb.WriteString(" " + parts[1])
		}
	}
	return sb.String()
}

func regsBlock(sys *system.System) string {
	regs := sys.Registers()
	var sb strings.Builder
	for i := range 16 {
		if i > 0 {
			sb.WriteByte('\n')
		}
		label := fmt.Sprintf("V%X", i)
		val := regValColor.Render(fmt.Sprintf("0x%02X", regs[i]))
		sb.WriteString(fmt.Sprintf("%-2s = %s (%3d)", label, val, regs[i]))
	}
	sb.WriteString("\n" + fmt.Sprintf("PC = %s", addrColor.Render(fmt.Sprintf("0x%03X", sys.PC()))))
	sb.WriteString("   " + fmt.Sprintf("I  = %s", addrColor.Render(fmt.Sprintf("0x%03X", sys.IReg()))))
	sb.WriteString("\n" + fmt.Sprintf("DT = %s", regValColor.Render(fmt.Sprintf("%-3d", sys.DelayTimer()))))
	sb.WriteString("   " + fmt.Sprintf("ST = %s", regValColor.Render(fmt.Sprintf("%-3d", sys.SoundTimer()))))
	if sys.AudioMuted() {
		sb.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("audio off"))
	} else {
		sb.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("audio on"))
	}
	return sb.String()
}

func keysBlock(sys *system.System) string {
	state := sys.KeyState()
	order := []byte{1, 2, 3, 0xC, 4, 5, 6, 0xD, 7, 8, 9, 0xE, 0xA, 0, 0xB, 0xF}
	labels := []string{"1", "2", "3", "C", "4", "5", "6", "D", "7", "8", "9", "E", "A", "0", "B", "F"}
	var sb strings.Builder
	for i, k := range order {
		if i > 0 && i%4 == 0 {
			sb.WriteByte('\n')
		} else if i%4 != 0 {
			sb.WriteString("  ")
		}
		if state[k] {
			sb.WriteString(keyOnStyle.Render(labels[i]))
		} else {
			sb.WriteString(keyOffStyle.Render(labels[i]))
		}
	}
	return sb.String()
}

func printDis(sys *system.System, pc uint16, n int) {
	fmt.Println()
	for i := range n {
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

func printRegs(sys *system.System) {
	regs := sys.Registers()
	fmt.Println()
	for i := range 16 {
		fmt.Printf("  V%X = 0x%02X (%3d)", i, regs[i], regs[i])
		if i%4 == 3 {
			fmt.Println()
		}
	}
	fmt.Printf("  PC = 0x%03X   I = 0x%03X\n\n", sys.PC(), sys.IReg())
}

func printMem(sys *system.System, addr uint16, n int) {
	mem := sys.MemorySlice(addr, n)
	fmt.Println()
	for i, b := range mem {
		if i%16 == 0 {
			fmt.Printf("  0x%03X: ", addr+uint16(i))
		}
		fmt.Printf("%02X ", b)
		if i%16 == 15 {
			fmt.Println()
		}
	}
	if len(mem)%16 != 0 {
		fmt.Println()
	}
	fmt.Println()
}

func printKeys(sys *system.System) {
	state := sys.KeyState()
	fmt.Println()
	fmt.Println("  1 2 3 C")
	fmt.Println("  4 5 6 D")
	fmt.Println("  7 8 9 E")
	fmt.Println("  A 0 B F")
	fmt.Println()
	order := []byte{1, 2, 3, 0xC, 4, 5, 6, 0xD, 7, 8, 9, 0xE, 0xA, 0, 0xB, 0xF}
	for i, k := range order {
		mark := "."
		if state[k] {
			mark = "#"
		}
		fmt.Printf("  %s", mark)
		if i%4 == 3 {
			fmt.Println()
		}
	}
	fmt.Println()
}
