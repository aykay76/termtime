package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"time"

	"golang.org/x/term"
)

var w, h int

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

func frame() {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println(err)
	}

	if width != w || height != h {
		w = width
		h = height
		redraw()
	}
}

func redraw() {
	clear()
	drawBox(5, 5, w-10, h-10)
	printCenter(5, " Hello, World! ")
}

func main() {
	// raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
	}

	// handle interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		term.Restore(int(os.Stdin.Fd()), oldState)
		showCursor()
		echoOn()
		rmcup()
		os.Exit(1)
	}()

	smcup()
	clear()
	hideCursor()
	echoOff()
	enableMouseTracking()

	// get the terminal size
	w, h, _ = term.GetSize(int(os.Stdout.Fd()))

	redraw()

	// how to save and restore screen for dialog boxes
	// smcup()
	// printCenterBox(3, "Welcome to some program!\r\n\r\nThis is some introduction text that will be displayed for 2 seconds. \r\n\r\nPress 'q' to quit at any time")
	// time.Sleep(2 * time.Second)
	// rmcup()

	go func() {
		for {
			move(1, 1)
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Split(bufio.ScanBytes)
			for scanner.Scan() {
				r := scanner.Bytes()
				switch r[0] {
				case 0x1b:
					scanner.Scan()
					r2 := scanner.Bytes()
					switch r2[0] {
					case 0x31:
						scanner.Scan()
						r3 := scanner.Bytes()
						switch r3[0] {
						case 0x35:
							scanner.Scan()
							r4 := scanner.Bytes()
							switch r4[0] {
							case 0x7e:
								fmt.Print("F5")
							default:
								fmt.Printf("%x", r4[0])
							}
						case 0x37:
							scanner.Scan()
							r4 := scanner.Bytes()
							switch r4[0] {
							case 0x7e:
								fmt.Print("F6")
							default:
								fmt.Printf("%x", r4[0])
							}
						default:
							fmt.Printf("%x", r3[0])
						}
					case 0x5b:
						scanner.Scan()
						r3 := scanner.Bytes()
						switch r3[0] {
						case 0x31:
							scanner.Scan()
							r4 := scanner.Bytes()
							switch r4[0] {
							case 0x35:
								scanner.Scan()
								r5 := scanner.Bytes()
								switch r5[0] {
								case 0x7e:
									fmt.Print("F5")
								default:
									fmt.Printf("35 %x", r5[0])
								}
							case 0x37:
								scanner.Scan()
								r5 := scanner.Bytes()
								switch r5[0] {
								case 0x7e:
									fmt.Print("F6")
								default:
									fmt.Printf("37 %x", r5[0])
								}
							case 0x38:
								scanner.Scan()
								r5 := scanner.Bytes()
								switch r5[0] {
								case 0x7e:
									fmt.Print("F7")
								default:
									fmt.Printf("38 %x", r5[0])
								}
							case 0x39:
								scanner.Scan()
								r5 := scanner.Bytes()
								switch r5[0] {
								case 0x7e:
									fmt.Print("F8")
								}
							default:
								fmt.Printf("31 %x", r4[0])
							}
						case 0x32:
							scanner.Scan()
							r4 := scanner.Bytes()
							switch r4[0] {
							case 0x30:
								scanner.Scan()
								r5 := scanner.Bytes()
								switch r5[0] {
								case 0x7e:
									fmt.Print("F9")
								default:
									fmt.Printf("30 %x", r5[0])
								}
							case 0x31:
								scanner.Scan()
								r5 := scanner.Bytes()
								switch r5[0] {
								case 0x7e:
									fmt.Print("F10")
								default:
									fmt.Printf("37 %x", r5[0])
								}
							case 0x32:
								scanner.Scan()
								r5 := scanner.Bytes()
								switch r5[0] {
								case 0x7e:
									fmt.Print("F11")
								default:
									fmt.Printf("38 %x", r5[0])
								}
							case 0x34:
								scanner.Scan()
								r5 := scanner.Bytes()
								switch r5[0] {
								case 0x7e:
									fmt.Print("F12")
								}
							default:
								fmt.Printf("32 %x", r4[0])
							}
						case 0x41:
							fmt.Print("up")
						case 0x42:
							fmt.Print("down")
						case 0x43:
							fmt.Print("right")
						case 0x44:
							fmt.Print("left")
						case 0x4d: // M
							move(1, 1)
							fmt.Print("mouse ")
							scanner.Scan()
							r4 := scanner.Bytes()
							switch r4[0] {
							case 0x20:
								fmt.Printf("left,")
							case 0x21:
								fmt.Printf("middle,")
							case 0x22:
								fmt.Printf("right,")
							case 0x23:
								fmt.Printf("up,")
							case 0x60:
								fmt.Printf("scroll up,")
							case 0x61:
								fmt.Printf("scroll down,")
							case 0x62:
								fmt.Printf("scroll left,")
							case 0x63:
								fmt.Printf("scroll right,")
							default:
								fmt.Printf("4d %x", r4[0])
							}
							scanner.Scan()
							r5 := scanner.Bytes()
							fmt.Printf("x=%d,", r5[0]-32)

							scanner.Scan()
							r6 := scanner.Bytes()
							fmt.Printf("y=%d", r6[0]-32)
						default:
							fmt.Printf("5b %x", r3[0])
						}
					case 0x4f:
						scanner.Scan()
						r3 := scanner.Bytes()
						switch r3[0] {
						case 0x48:
							fmt.Print("home")
						case 0x46:
							fmt.Print("end")
						case 0x50:
							fmt.Print("F1")
						case 0x51:
							fmt.Print("F2")
						case 0x52:
							fmt.Print("F3")
						case 0x53:
							fmt.Print("F4")
						default:
							fmt.Printf("4f %x", r3[0])
						}
					default:
						fmt.Printf("1b %x", r2[0])
					}
				case 'q':
					c <- os.Interrupt
				}
			}
		}
	}()

	// loop forever checking for keyboard and mouse input
	for {
		time.Sleep(40 * time.Millisecond)
		frame()
	}
}
