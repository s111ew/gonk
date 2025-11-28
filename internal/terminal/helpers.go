package terminal

func IsCtrl(c byte) bool {
	return c <= 31 || c == 127
}

func CtrlKey(c byte) byte {
	return c & 0x1F
}
