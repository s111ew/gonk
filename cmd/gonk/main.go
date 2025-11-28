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
		c := buf[0]
		// exit with "q"
		if err != nil || c == terminal.CtrlKey('q') {
			break
		}

		// print new line for enter
		if c == 13 {
			fmt.Print("\r\n")
		} else if terminal.IsCtrl(c) {
			fmt.Printf("%d\r\n", c)
		} else {
			fmt.Printf("%d (%c)", c, c)
		}
	}
}
