Output demo.gif

Set Shell "fish"

Set Theme "GruvboxDark"

Set FontSize 16
Set FontFamily "FiraCode Nerd Font"

Set Width 1200
Set Height 700

Set Padding 40
Set Margin 40
Set MarginFill "#d79921"
Set BorderRadius 10

Set TypingSpeed 80ms

Set Framerate 60

Hide
Type "set -g fish_autosuggestion_enabled 0"
Enter
Type "set fish_color_command 458588 --bold"
Enter
Type "set fish_color_param ebdbb2"
Enter
Type "set fish_color_quote b8bb26"
Enter
Type "set fish_color_error cc241d --bold"
Enter
Sleep 500ms
Show

Ctrl+L
Sleep 1s

# Scene 1: launch the debugger
Type "./g8dbg test_roms/IBM_logo.ch8"
Sleep 1s
Enter
Sleep 2s

# Scene 2: show help
Type "h"
Sleep 1s
Enter
Sleep 5s

# Scene 3: disassemble
Type "dis 8"
Sleep 1s
Enter
Sleep 4s

# Scene 4: step a few instructions
Type "step 3"
Sleep 1s
Enter
Sleep 4s

# Scene 5: inspect registers
Type "regs"
Sleep 1s
Enter
Sleep 4s

# Scene 6: dump memory at ROM start
Type "m 0x200 32"
Sleep 1s
Enter
Sleep 4s

# Scene 7: set a breakpoint and continue
Type "b 0x228"
Sleep 1s
Enter
Sleep 2s
Type "c"
Sleep 1s
Enter
Sleep 4s

# Scene 8: inspect state at breakpoint
Type "regs"
Sleep 1s
Enter
Sleep 4s

# Scene 9: timers
Type "timers"
Sleep 1s
Enter
Sleep 4s

# Scene 10: quit
Type "q"
Sleep 1s
Enter
Sleep 4s
