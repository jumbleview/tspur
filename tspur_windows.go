//go:build windows
// +build windows

package main

import (
	"fmt"
	"os/exec"
)

// SetDimensions sets the size of the console (alas it  does not work for Windows 11 console)
func SetDimensions(cols int, lines int) {
	columns := fmt.Sprintf("cols=%d", cols)
	rows := fmt.Sprintf("lines=%d", lines)
	exec.Command("mode", "con:", columns, rows).Output()
}
