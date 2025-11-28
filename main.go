package main

import (
	"fmt"
	"os"

	terminalconfig "github.com/s111ew/gonk/terminal_config"
)

func main() {
	err := terminalconfig.EnableRawMode()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer func() {
		if err := terminalconfig.DisableRawMode(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to restore terminal:", err)
		}
	}()

	var buf [1]byte
	for {
		n, err := os.Stdin.Read(buf[:])
		if err != nil || n == 0 {
			break
		}
	}
}
