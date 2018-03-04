// up is the ultimate pipe composer/editor. It helps building Linux pipelines
// in a terminal-based UI interactively, with live preview of command results.
package main

import (
	"fmt"
	"io"
	"os"

	termbox "github.com/akavel/termbox-go"
	"github.com/mattn/go-isatty"
)

func main() {
	// TODO: Without below block, we'd hang with no piped input (see github.com/peco/peco, mattn/gof, fzf, etc.)
	if isatty.IsTerminal(os.Stdin.Fd()) {
		fmt.Fprintln(os.Stderr, "error: up requires some data piped on standard input, e.g.: `echo hello world | up`")
		os.Exit(1)
	}

	// Init TUI code
	// TODO: maybe try gocui or tcell?
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// In background, start collecting input from stdin to internal buffer of size 40 MB, then pause it
	go collect()

	var (
		prompt  = []rune("| ")
		command = []rune{}
		cursor  = 0
	)

	// Main loop
main_loop:
	for {
		// Draw command input line
		for x, ch := range prompt {
			termbox.SetCell(x, 0, ch, termbox.ColorWhite, termbox.ColorBlue)
		}
		for x, ch := range command {
			termbox.SetCell(x+len(prompt), 0, ch, termbox.ColorWhite, termbox.ColorBlue)
		}
		termbox.SetCursor(len(prompt)+cursor, 0)
		termbox.Flush()

		// Handle events
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Ch != 0 {
				// insert key into command (https://github.com/golang/go/wiki/SliceTricks#insert)
				command = append(command, 0)
				copy(command[cursor+1:], command[cursor:])
				command[cursor] = ev.Ch
				cursor++
				continue main_loop
			}
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				// quit
				return
			// handle command-line editing keys
			default:
			}
		}
	}

	// TODO: using tcell, edit a command in bash format in multiline input box (or jroimartin/gocui?)
	//       NOTE: gocui has trouble if we capture stdin. Try butchering ("total modding") peco/peco instead.
	// TODO: run it automatically in bg after first " " (or ^Enter), via `bash -c`
	// TODO: auto-kill the child process on any edit
	// TODO: allow scrolling the output preview with pgup/pgdn keys
	// TODO: [LATER] Ctrl-O shows input via `less` or $PAGER
	// TODO: ^X - save into executable file upN.sh (with #!/bin/bash) and quit
	// TODO: [LATER] allow increasing size of input buffer with some key
	// TODO: [LATER] on ^X, leave TUI and run the command through buffered input, then unpause rest of input
	// TODO: [LATER] allow adding more elements of pipeline (initially, just writing `foo | bar` should work)
	// TODO: [LATER] allow invocation with partial command, like: `up grep -i`
	// TODO: [LATER][MAYBE] allow reading upN.sh scripts
	// TODO: [LATER] auto-save and/or save on Ctrl-S or something
	// TODO: [MUCH LATER] readline-like rich editing support? and completion?
	// TODO: [MUCH LATER] integration with fzf? and pindexis/marker?
	// TODO: [LATER] forking and unforking pipelines
	// TODO: [LATER] capture output of a running process (see: https://stackoverflow.com/q/19584825/98528)
	// TODO: [LATER] richer TUI:
	// - show # of read lines & kbytes
	// - show status (errorlevel) of process, or that it's still running (also with background colors)
	// - allow copying and pasting to/from command line
	// TODO: [LATER] allow connecting external editor (become server/engine via e.g. socket)
	// TODO: [LATER] become pluggable into http://luna-lang.org
	// TODO: [LATER][MAYBE] allow "plugins" ("combos" - commands with default options) e.g. for Lua `lua -e`+auto-quote, etc.
	// TODO: [LATER] make it more friendly to infrequent Linux users by providing "descriptive" commands like "search" etc.
	// TODO: [LATER] advertise on: HN, r/programming, r/golang, r/commandline, r/linux; data exploration? data science?
}

func collect() {
	const bufsize = 40 * 1024 * 1024 // 40 MB
	buf := make([]byte, bufsize)
	// TODO: read gradually what is available and show progress
	n, err := io.ReadFull(os.Stdin, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		panic(err)
	}
	buf = buf[:n]
	// TODO: use buf somewhere
}
