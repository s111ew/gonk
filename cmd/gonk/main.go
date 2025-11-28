package main

import (
	"fmt"
	"os"

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

	var buf [1]byte

	for {
		_, err := os.Stdin.Read(buf[:])
		// exit with "q"
		if err != nil || buf[0] == 113 {
			break
		}

		// print new line for enter
		if buf[0] == 13 {
			fmt.Print("\r\n")
		} else {
			fmt.Print(string(buf[0]))
		}
	}
}
