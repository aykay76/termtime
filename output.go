package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func clear() {
	fmt.Print("\033[2J")
}

func enableMouseTracking() {
	fmt.Print("\033[?1000h")
}

func disableMouseTracking() {
	fmt.Print("\033[?1000l")
}

func hideCursor() {
	fmt.Print("\033[?25l")
}

func echoOff() {
	fmt.Print("\033[?7l")
}

func smcup() {
	fmt.Print("\033[?1049h")
}

func rmcup() {
	fmt.Print("\033[?1049l")
}

func echoOn() {
	fmt.Print("\033[?7h")
}

func showCursor() {
	fmt.Print("\033[?25h")
}

func move(x, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

func printAt(x, y int, s string) {
	move(x, y)
	fmt.Print(s)
}

func printCenter(y int, s string) {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	x := (width - len(s)) / 2
	printAt(x, y, s)
}

func printCenterln(y int, s string) {
	printCenter(y, s)
	fmt.Println()
}

func printCenterlnf(y int, format string, a ...interface{}) {
	printCenterln(y, fmt.Sprintf(format, a...))
}

func printCenterf(y int, format string, a ...interface{}) {
	printCenter(y, fmt.Sprintf(format, a...))
}

func printCenterBox(y int, s string) {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	x := (width - len(s)) / 2
	drawBox(x-1, y-1, len(s)+2, 3)
	printAt(x, y, s)
}

func printCenterBoxf(y int, format string, a ...interface{}) {
	printCenterBox(y, fmt.Sprintf(format, a...))
}

func printCenterBoxln(y int, s string) {
	printCenterBox(y, s)
	fmt.Println()
}

func drawBox(x, y, width, height int) {
	move(x, y)

	// draw a box
	fmt.Print("\033(0")
	fmt.Print("l")
	for i := 0; i < width-2; i++ {
		fmt.Print("q")
	}
	fmt.Print("k")
	y++
	for i := 0; i < height-2; i++ {
		move(x, y)

		fmt.Print("x")
		fmt.Printf("\033[%dC", width-2)
		fmt.Print("x")
		y++
	}

	move(x, y)
	fmt.Print("m")
	for i := 0; i < width-2; i++ {
		fmt.Print("q")
	}
	fmt.Print("j")
	fmt.Print("\033(B")
}

func setBackground(color int) {
	fmt.Printf("\033[48;5;%dm", color)
}
