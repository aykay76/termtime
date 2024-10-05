package main

import (
	"fmt"
	"testing"
)

func TestDrawBox(t *testing.T) {
	// var buf bytes.Buffer
	// // TODO: see this for blocking consumer to output directly to screen - can use this to redirect stdout to a buffer
	// old := os.Stdout
	// r, w, _ := os.Pipe()
	// os.Stdout = w

	fmt.Print("\033(0")
	fmt.Print("abcdefghijklmnopqrstuvwxyz1234567890")
	fmt.Print("\033(B")

	// w.Close()
	// os.Stdout = old
	// buf.ReadFrom(r)

	// expectedOutput := "\033[1;1H\033(0lqqqqqqqqk\033[2;1Hx        x\033[3;1Hx        x\033[4;1Hx        x\033[5;1Hmqqqqqqqqj\033(B"
	// if buf.String() != expectedOutput {
	// 	t.Errorf("Expected %q but got %q", expectedOutput, buf.String())
	// }
}
