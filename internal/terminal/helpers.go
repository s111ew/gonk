package terminal

// returns true if character is 'ctrl' modified
func IsCtrl(c byte) bool {
	return c <= 31 || c == 127
}

// return char code for the 'ctrl' modified version of
// given key
func CtrlKey(c byte) byte {
	return c & 0x1F
}
