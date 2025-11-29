package editor

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/s111ew/gonk/internal/terminal"
)

// package level unique error for signalling to main
// function that the program should quit
var ErrQuit = errors.New("quit")

// populate the user terminal config struct with window dimensions
func InitEditor() error {
	if err := terminal.GetWindowSize(&terminal.Config); err != nil {
		return err
	}
	return nil
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
	var buf strings.Builder

	buf.WriteString("\x1b[2J")
	buf.WriteString("\x1b[H")

	drawRows(&buf)

	buf.WriteString("\x1b[H")
	os.Stdout.Write([]byte(buf.String()))
}

// draw rows on user terminal
func drawRows(buf *strings.Builder) {
	for y := 0; y < terminal.Config.ScreenRows; y++ {
		buf.WriteString("~")

		if y < terminal.Config.ScreenRows-1 {
			buf.WriteString("\r\n")
		}
	}
}
