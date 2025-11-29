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

const (
	ARROW_LEFT  = 255
	ARROW_RIGHT = 254
	ARROW_UP    = 253
	ARROW_DOWN  = 252
)

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

		if n == 0 {
			continue
		}

		if buf[0] == '\x1b' {
			var seq [3]byte

			n1, err := os.Stdin.Read(seq[:1])
			if err != nil {
				return 0, err
			}
			if n1 != 1 {
				return '\x1b', nil
			}

			n2, err := os.Stdin.Read(seq[1:2])
			if err != nil {
				return 0, err
			}
			if n2 != 1 {
				return '\x1b', nil
			}

			if seq[0] == '[' {
				switch seq[1] {
				case 'A':
					return ARROW_UP, nil
				case 'B':
					return ARROW_DOWN, nil
				case 'C':
					return ARROW_RIGHT, nil
				case 'D':
					return ARROW_LEFT, nil
				}
			}

			return '\x1b', nil
		}

		return buf[0], nil
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
	case ARROW_UP:
	case ARROW_DOWN:
	case ARROW_LEFT:
	case ARROW_RIGHT:
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
	case ARROW_LEFT:
		terminal.Config.CursorX--
	case ARROW_RIGHT:
		terminal.Config.CursorX++
	case ARROW_UP:
		terminal.Config.CursorY--
	case ARROW_DOWN:
		terminal.Config.CursorY++
	}
}
