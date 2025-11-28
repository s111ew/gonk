package main

import "os"

func main() {
	var buf [1]byte
	for {
		n, err := os.Stdin.Read(buf[:])
		if err != nil || n == 0 {
			break
		}
	}
}
