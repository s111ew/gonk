package terminal

import (
	"os"

	"golang.org/x/sys/unix"
)

// save a reference to the user's terminal settings
// so we can restore them when we exit the program
var origTermios *unix.Termios

// disable ECHO mode in the terminal
// (user can't see what is typed in the terminal)
func EnableRawMode() error {
	// grab a reference to the os' standard input stream
	// (where we read the terminal input from)
	fd := int(os.Stdin.Fd())

	// fetch current terminal settings
	termios, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	if err != nil {
		return err
	}

	origTermios = termios

	// clear the setting for 'terminal echo' and 'canonical'
	raw := *termios
	raw.Lflag &^= unix.ECHO | unix.ICANON

	// flush any pending input before applying changes
	err = unix.IoctlSetTermios(fd, unix.TIOCSETAF, &raw)
	if err != nil {
		return err
	}

	return nil
}

// re-enable ECHO mode in the terminal
func DisableRawMode() error {
	fd := int(os.Stdin.Fd())
	err := unix.IoctlSetTermios(fd, unix.TIOCSETAF, origTermios)
	return err
}
