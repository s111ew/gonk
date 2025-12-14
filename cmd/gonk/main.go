package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/s111ew/gonk/internal/editor"
	"github.com/s111ew/gonk/internal/terminal"
)

func main() {
	// enable 'raw' mode in user's terminal
	if err := terminal.EnableRawMode(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to activate no-echo mode:", err)
		return
	}

	// intialise user editor configuration
	if err := editor.InitEditor(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialise editor:", err)
		return
	}

	editor.EditorOpen(os.Args[1])

	// disable 'raw' mode in user's terminal when program exits
	defer func() {
		if err := terminal.DisableRawMode(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to deactivate no-echo mode:", err)
		}
	}()

	for {
		// listen for key press
		editor.RefreshScreen()
		err := editor.ProcessKeyPress()
		if err != nil {
			editor.RefreshScreen()
			// check for ctrl(q) press and exit program
			if errors.Is(err, editor.ErrQuit) {
				fmt.Println("exiting...")
				return
			}
			log.Fatal(err)
		}
	}
}
