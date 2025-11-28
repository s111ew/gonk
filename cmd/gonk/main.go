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
	err := terminal.EnableRawMode()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to activate no-echo mode:", err)
		return
	}

	defer func() {
		if err := terminal.DisableRawMode(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to deactivate no-echo mode:", err)
		}
	}()

	for {
		err := editor.ProcessKeyPress()
		if err != nil {
			if errors.Is(err, editor.ErrQuit) {
				fmt.Println("exiting...")
				return
			}
			log.Fatal(err)
		}
	}
}
