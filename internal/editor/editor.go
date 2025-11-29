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
			// escape sequence
			var seq []byte
			seq = append(seq, buf[0])

			// read the next two bytes minimum
			for i := 0; i < 2; i++ {
				_, err := os.Stdin.Read(buf[:])
				if err != nil {
					return 0, err
				}
				seq = append(seq, buf[0])
			}

			// read until letter
			for {
				b := seq[len(seq)-1]
				if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
					break
				}
				_, err := os.Stdin.Read(buf[:])
				if err != nil {
					return 0, err
				}
				seq = append(seq, buf[0])
			}

			end := seq[len(seq)-1]

			switch end {
			case 'A':
				// log.Fatalf("RETURNING: %v'\n", ARROW_UP)
				return ARROW_UP, nil
			case 'B':
				// log.Fatalf("RETURNING: %v'\n", ARROW_DOWN)
				return ARROW_DOWN, nil
			case 'C':
				// log.Fatalf("RETURNING: %v'\n", ARROW_RIGHT)
				return ARROW_RIGHT, nil
			case 'D':
				// log.Fatalf("RETURNING: %v'\n", ARROW_LEFT)
				return ARROW_LEFT, nil
			}

			return 0, nil
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
	// quit out if input is ctrl + c
	case terminal.CtrlKey('q'):
		return ErrQuit
	// move cursor around with 'wasd'
	case ARROW_UP:
		moveCursor(c)
	case ARROW_DOWN:
		moveCursor(c)
	case ARROW_LEFT:
		moveCursor(c)
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

	buf.WriteString(fmt.Sprintf("\x1b[%d;%dH", terminal.Config.CursorY+1, terminal.Config.CursorX+1))
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
		if terminal.Config.CursorX > 0 {
			terminal.Config.CursorX--
		}
	case ARROW_RIGHT:
		if terminal.Config.CursorX < terminal.Config.ScreenCols-1 {
			terminal.Config.CursorX++
		}
	case ARROW_UP:
		if terminal.Config.CursorY > 0 {
			terminal.Config.CursorY--
		}
	case ARROW_DOWN:
		if terminal.Config.CursorY < terminal.Config.ScreenRows-1 {
			terminal.Config.CursorY++
		}
	}
}
