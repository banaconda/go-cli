package libutil

import "fmt"

func ClearLine() {
	fmt.Printf("\033[2K\r")
}

func ClearScreen() {
	fmt.Printf("\033[2J")
}

// move cursor
func MoveCursor(direction int) {
	if direction == 0 {
		return
	}

	if direction > 0 {
		fmt.Printf("\033[%dC", direction)
	} else {
		fmt.Printf("\033[%dD", -direction)
	}
}
