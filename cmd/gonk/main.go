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
	err := terminal.EnableRawMode()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to activate no-echo mode:", err)
		return
	}

	// disable 'raw' mode in user's terminal when program exits
	defer func() {
		if err := terminal.DisableRawMode(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to deactivate no-echo mode:", err)
		}
	}()

	for {
		// listen for key press
		err := editor.ProcessKeyPress()
		if err != nil {
			// check for ctrl(q) press and exit program
			if errors.Is(err, editor.ErrQuit) {
				fmt.Println("exiting...")
				return
			}
			log.Fatal(err)
		}
	}
}
