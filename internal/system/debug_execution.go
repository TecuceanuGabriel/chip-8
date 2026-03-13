package system

import (
	"fmt"
	"os"
	"time"
)

// runDebugFrame runs one frame's worth of execution in debug mode.
// When paused it blocks until the debugger sends a command; timers only
// tick when instructions actually execute.
func (system *System) runDebugFrame() {
	if system.isPaused {
		system.waitForDebugCmd()
		return
	}

	executedAny := false
	for range nrInstPerFrame {
		select {
		case cmd := <-system.debug.debugChan:
			system.handleDebugCmd(cmd)
			if system.isPaused {
				system.waitForDebugCmd()
				return
			}
		default:
		}

		if system.debug.breakpoints[system.pc] {
			system.isPaused = true
			system.debug.eventChan <- EventBreakpoint{PC: system.pc}
			system.waitForDebugCmd()
			return
		}

		system.execInstruction()
		executedAny = true

		if system.debug.stepMode {
			system.debug.stepsLeft--
			if system.debug.stepsLeft <= 0 {
				system.debug.stepMode = false
				system.isPaused = true
				system.debug.eventChan <- EventStep{PC: system.pc}
				system.waitForDebugCmd()
				return
			}
		}
	}

	if executedAny {
		system.updateTimers()
	}
}

// waitForDebugCmd blocks until the debugger unpauses execution, refreshing
// the window on every frame tick so the OS does not consider it frozen.
func (system *System) waitForDebugCmd() {
	refresh := time.NewTicker(time.Second / targetFPS)
	defer refresh.Stop()
	win := system.display.GetWindow()

	for system.isPaused {
		select {
		case cmd := <-system.debug.debugChan:
			system.handleDebugCmd(cmd)
		case <-refresh.C:
			win.Update()
		}
	}
}

// handleDebugCmd applies a single command from the debugger goroutine.
func (system *System) handleDebugCmd(cmd DebugCmd) {
	switch c := cmd.(type) {
	case CmdStep:
		system.debug.stepMode = true
		system.debug.stepsLeft = c.Count
		system.isPaused = false
	case CmdContinue:
		system.debug.stepMode = false
		system.isPaused = false
	case CmdSetBreakpoint:
		system.debug.breakpoints[c.Addr] = true
	case CmdRemoveBreakpoint:
		delete(system.debug.breakpoints, c.Addr)
	case CmdSetKey:
		system.keyState[c.Key] = true
	case CmdReleaseKey:
		system.keyState[c.Key] = false
	case CmdReset:
		system.reset()
		system.debug.stepMode = false
		system.debug.stepsLeft = 0
		system.debug.eventChan <- EventStep{PC: system.pc}
	case CmdQuit:
		os.Exit(0)
	default:
		fmt.Printf("unknown debug command: %T\n", cmd)
	}
}
