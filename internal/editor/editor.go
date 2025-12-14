package editor

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/s111ew/gonk/internal/terminal"
)

const (
	WELCOME_MSG = "Gonk -- "
	VERSION     = "0.0.1"

	ARROW_LEFT  = 255
	ARROW_RIGHT = 254
	ARROW_UP    = 253
	ARROW_DOWN  = 252
	PAGE_UP     = 251
	PAGE_DOWN   = 250
	HOME_KEY    = 249
	END_KEY     = 248
	DEL_KEY     = 247
)

// package level unique error for signalling to main
// function that the program should quit
var ErrQuit = errors.New("quit")

// populate the user terminal config struct with window dimensions
func InitEditor() error {
	terminal.Config.CursorX = 0
	terminal.Config.CursorY = 0
	terminal.Config.NumRows = 0

	if err := terminal.GetWindowSize(&terminal.Config); err != nil {
		return err
	}
	return nil
}

func EditorOpen(filename string) {
	// Open the file
	fp, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("fopen: %v", err))
	}
	defer fp.Close()

	// Read the first line
	scanner := bufio.NewScanner(fp)
	if scanner.Scan() {
		line := scanner.Text()

		// Trim newline and carriage return characters
		line = strings.TrimRight(line, "\r\n")

		// Assign to the editor structure
		terminal.Config.Row.Size = len(line)
		terminal.Config.Row.Text = line
		terminal.Config.NumRows = 1
	}

	// Check for scanner error
	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("getline: %v", err))
	}

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

		// normal single byte key press
		if buf[0] != '\x1b' {
			return buf[0], nil
		}

		// ESC sequence starts
		seq := []byte{'\x1b'}

		// read next byte (must exist for an escape sequence)
		_, err = os.Stdin.Read(buf[:])
		if err != nil {
			return 0, err
		}
		seq = append(seq, buf[0])

		// if not CSI, ignore
		if seq[1] == '[' {
			// read next byte
			_, err = os.Stdin.Read(buf[:])
			if err != nil {
				return 0, err
			}
			seq = append(seq, buf[0])

			// if it's a digit, expect something like "ESC [ 5 ~"
			if seq[2] >= '0' && seq[2] <= '9' {
				// read until '~'
				for {
					_, err = os.Stdin.Read(buf[:])
					if err != nil {
						return 0, err
					}
					seq = append(seq, buf[0])
					if buf[0] == '~' {
						break
					}
				}

				switch seq[2] {
				case '1':
					return HOME_KEY, nil
				case '3':
					return DEL_KEY, nil
				case '4':
					return END_KEY, nil
				case '5':
					return PAGE_UP, nil
				case '6':
					return PAGE_DOWN, nil
				case '7':
					return HOME_KEY, nil
				case '8':
					return END_KEY, nil
				}
				return 0, nil
			}
		} else if seq[1] == 'O' {
			switch seq[2] {
			case 'H':
				return HOME_KEY, nil
			case 'F':
				return END_KEY, nil
			}
		}

		// arrow keys
		last := seq[2]

		switch last {
		case 'A':
			return ARROW_UP, nil
		case 'B':
			return ARROW_DOWN, nil
		case 'C':
			return ARROW_RIGHT, nil
		case 'D':
			return ARROW_LEFT, nil
		case 'H':
			return HOME_KEY, nil
		case 'F':
			return END_KEY, nil
		}

		return 0, nil
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

	case PAGE_UP, PAGE_DOWN:
		// how many times to move up/down
		n := terminal.Config.ScreenRows
		// move up/down until at top
		for n > 0 {
			if c == PAGE_UP {
				moveCursor(ARROW_UP)
			} else {
				moveCursor(ARROW_DOWN)
			}
			n--
		}

	// move cursor around with arrow keys
	case ARROW_UP, ARROW_DOWN, ARROW_LEFT, ARROW_RIGHT:
		moveCursor(c)

	case HOME_KEY:
		terminal.Config.CursorX = 0

	case END_KEY:
		terminal.Config.CursorX = terminal.Config.ScreenCols - 1
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

	for y := 0; y < terminal.Config.ScreenRows; y++ {
		if terminal.Config.NumRows == 0 && y == 0 {
			msg := WELCOME_MSG + VERSION

			padding := (terminal.Config.ScreenCols - len(msg)) / 2
			if padding > 0 {
				for range padding {
					buf.WriteString(" ")
				}
			}
			buf.WriteString(msg)
			buf.WriteString("\r\n")
		}

		// if line is below the text we have, then print empty space
		if y >= terminal.Config.NumRows {
			buf.WriteString("~")
			buf.WriteString("\x1b[K")

			if y < terminal.Config.ScreenRows-1 {
				buf.WriteString("\r\n")
			}

			// else print out our text row
		} else {
			length := terminal.Config.Row.Size
			row := terminal.Config.Row.Text
			if length > terminal.Config.ScreenCols {
				row = truncateString(row, terminal.Config.ScreenCols)
			}
			buf.WriteString(row)
			buf.WriteString("\x1b[K")
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
