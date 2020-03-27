// +build windows

package main

import (
	"fmt"
	"os/exec"
)

// SetDimension sets the size of console
func SetDimension(cols int, lines int) {
	columns := fmt.Sprintf("cols=%d", cols)
	rows := fmt.Sprintf("lines=%d", lines)
	setMode := exec.Command("mode", "con:", columns, rows)
	setMode.Run()
}
