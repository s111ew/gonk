package editor

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/s111ew/gonk/internal/terminal"
)

var WELCOME_MSG string = "Gonk -- "
var VERSION string = "0.0.1"

// package level unique error for signalling to main
// function that the program should quit
var ErrQuit = errors.New("quit")

// populate the user terminal config struct with window dimensions
func InitEditor() error {
	terminal.Config.CursorX = 0
	terminal.Config.CursorY = 0

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

	// evaluate key input
	switch c {
	// quit out of input is ctrl + c
	case terminal.CtrlKey('q'):
		return ErrQuit
	// move cursor around with 'wasd'
	case 'w':
	case 's':
	case 'a':
	case 'd':
		moveCursor(c)
	}

	return nil
}

// refresh user terminal screen
func RefreshScreen() {
	var buf strings.Builder

	buf.WriteString("\x1b[?25l")
	buf.WriteString("\x1b[H")

	drawRows(&buf)

	buf.WriteString(fmt.Sprintf("\x1b[%d;%dH", terminal.Config.CursorX+1, terminal.Config.CursorY+1))
	buf.WriteString("\x1b[?25h")

	os.Stdout.Write([]byte(buf.String()))
}

// draw rows on user terminal
func drawRows(buf *strings.Builder) {
	msg := WELCOME_MSG + VERSION
	for y := 0; y < terminal.Config.ScreenRows; y++ {
		buf.WriteString("~")
		if y == 0 {
			padding := ((terminal.Config.ScreenCols - len(msg)) / 2) - 1
			if padding > 0 {
				for range padding {
					buf.WriteString(" ")
				}
			}
			buf.WriteString(msg)
		}

		buf.WriteString("\x1b[K")
		if y < terminal.Config.ScreenRows-1 {
			buf.WriteString("\r\n")
		}
	}
}

func moveCursor(key byte) {
	switch key {
	case 'a':
		terminal.Config.CursorX--
	case 'd':
		terminal.Config.CursorX++
	case 'w':
		terminal.Config.CursorY--
	case 's':
		terminal.Config.CursorY++
	}
}
