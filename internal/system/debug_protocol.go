package system

// DebugCmd is the interface implemented by all commands sent from the debugger
// goroutine to the game loop.
type DebugCmd interface{ debugCmd() }

// CmdStep asks the game loop to execute Count instructions then pause.
type CmdStep struct{ Count int }

// CmdContinue asks the game loop to run freely until a breakpoint is hit.
type CmdContinue struct{}

// CmdSetBreakpoint registers a breakpoint at the given address.
type CmdSetBreakpoint struct{ Addr uint16 }

// CmdRemoveBreakpoint removes the breakpoint at the given address.
type CmdRemoveBreakpoint struct{ Addr uint16 }

// CmdSetKey injects a key-press into the emulator's key state.
type CmdSetKey struct{ Key byte }

// CmdReleaseKey injects a key-release into the emulator's key state.
type CmdReleaseKey struct{ Key byte }

// CmdQuit signals the game loop to exit.
type CmdQuit struct{}

func (CmdStep) debugCmd()             {}
func (CmdContinue) debugCmd()         {}
func (CmdSetBreakpoint) debugCmd()    {}
func (CmdRemoveBreakpoint) debugCmd() {}
func (CmdSetKey) debugCmd()           {}
func (CmdReleaseKey) debugCmd()       {}
func (CmdQuit) debugCmd()             {}

// DebugEvent is the interface implemented by all events sent from the game loop
// to the debugger goroutine.
type DebugEvent interface{ debugEvent() }

// EventBreakpoint is fired when execution reaches a breakpointed address.
type EventBreakpoint struct{ PC uint16 }

// EventStep is fired after the requested number of steps complete.
type EventStep struct{ PC uint16 }

func (EventBreakpoint) debugEvent() {}
func (EventStep) debugEvent()       {}

// debugSession holds breakpoints, step counters, and the channels used to
// communicate between the debugger goroutine and the game loop.
type debugSession struct {
	breakpoints map[uint16]bool
	stepMode    bool
	stepsLeft   int
	debugChan   chan DebugCmd
	eventChan   chan DebugEvent
}
