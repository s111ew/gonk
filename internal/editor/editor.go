package editor

import (
	"errors"
	"fmt"
	"os"

	"github.com/s111ew/gonk/internal/terminal"
)

var ErrQuit = errors.New("quit")

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
