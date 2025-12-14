package terminal

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// save a reference to the user's terminal settings
// so we can restore them when we exit the program
type TermConfig struct {
	OrigTermios *unix.Termios
	CursorX     int
	CursorY     int
	ScreenRows  int
	ScreenCols  int
	NumRows     int
	Row         Erow
}

type Erow struct {
	Size int
	Text string
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

var Config TermConfig

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

	Config.OrigTermios = termios

	// clear settings for 'echo', 'canonical' and treat ctrl
	// chars as regular char inputs
	raw := Config.OrigTermios
	raw.Iflag &^= unix.BRKINT | unix.ICRNL | unix.INPCK | unix.ISTRIP | unix.IXON
	raw.Oflag &^= unix.OPOST
	raw.Cflag |= unix.CS8
	raw.Lflag &^= unix.ECHO | unix.ICANON | unix.IEXTEN | unix.ISIG
	raw.Cc[unix.VMIN] = 1
	raw.Cc[unix.VTIME] = 0

	// flush any pending input before applying changes
	err = unix.IoctlSetTermios(fd, unix.TIOCSETAF, raw)
	if err != nil {
		return err
	}

	return nil
}

// re-enable ECHO mode in the terminal
func DisableRawMode() error {
	fd := int(os.Stdin.Fd())
	err := unix.IoctlSetTermios(fd, unix.TIOCSETAF, Config.OrigTermios)
	return err
}

// fetch the user's terminal size (rows , cols) and add to the
// global config struct
func GetWindowSize(config *TermConfig) error {
	ws := &winsize{}

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	if errno == 0 && ws.Col != 0 {
		config.ScreenRows = int(ws.Row)
		config.ScreenCols = int(ws.Col)
	} else {
		os.Stdout.Write([]byte("\x1b[999C\x1b[999B"))
		if err := getCursorPosition(config); err != nil {
			return err
		}
	}

	return nil
}

func getCursorPosition(config *TermConfig) error {
	// read response: ESC[n;mR
	buf := make([]byte, 32)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return err
	}

	// parse escape sequence
	// expecting: "\x1b[24;80R"
	re := regexp.MustCompile(`\x1b\[(\d+);(\d+)R`)
	m := re.FindStringSubmatch(string(buf[:n]))
	if len(m) != 3 {
		return fmt.Errorf("could not parse cursor position")
	}

	rows, _ := strconv.Atoi(m[1])
	cols, _ := strconv.Atoi(m[2])
	config.ScreenRows = rows
	config.ScreenCols = cols

	return nil
}
