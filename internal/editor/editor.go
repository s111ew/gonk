package editor

import (
	"errors"
	"fmt"
	"os"

	"github.com/s111ew/gonk/internal/terminal"
)

// package level unique error for signalling to main
// function that the program should quit
var ErrQuit = errors.New("quit")

func InitEditor() {
	terminal.GetWindowSize(&terminal.Config)
}

// wait for user key press and return it
func ReadKey() (byte, error) {
	var buf [1]byte

	for {
		n, err := os.Stdin.Read(buf[:])
		if err != nil {
			return 0, err
		}

		if n == 1 {
			return buf[0], nil
		}
	}
}

// evaluate pressed key
func ProcessKeyPress() error {
	c, err := ReadKey()
	if err != nil {
		return err
	}

	// exit with "q"
	if c == terminal.CtrlKey('q') {
		return ErrQuit
	} else if c == 13 {
		fmt.Print("\r\n")
	} else if terminal.IsCtrl(c) {
		fmt.Printf("%d\r\n", c)
	} else {
		fmt.Printf("%d (%c)", c, c)
	}

	return nil
}

// refresh user terminal screen
func RefreshScreen() {
	os.Stdout.Write([]byte("\x1b[2J"))
	os.Stdout.Write([]byte("\x1b[H"))

	drawRows()

	os.Stdout.Write([]byte("\x1b[H"))
}

// draw rows on user terminal
func drawRows() {
	// arbitrary 24 rows
	for y := 0; y < terminal.Config.ScreenRows; y++ {
		os.Stdout.Write([]byte("~\r\n"))
	}
}
