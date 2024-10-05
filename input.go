package main

import (
	"bufio"
	"fmt"
	"os"
)

func input(c chan os.Signal) {
	for {
		c0 := 0
		esc := 1
		csi := 2
		mouse := 3
		ss3 := 4
		state := c0
		parms := make([]int, 0)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanBytes)
		for scanner.Scan() {
			r := scanner.Bytes()[0]
			switch state {
			case c0:
				switch r {
				case 0x7:
					fmt.Print("BEL ")
				case 0x8:
					fmt.Print("BS ")
				case 0x9:
					fmt.Print("TAB ")
				case 0xa:
					fmt.Print("LF ")
				case 0xb:
					fmt.Print("VT ")
				case 0xc:
					fmt.Print("FF ")
				case 0xd:
					fmt.Print("CR ")
				case 0x1b:
					move(1, 1)
					fmt.Print("ESC ")
					state = esc
					parms = make([]int, 0)
				case 'q':
					c <- os.Interrupt
				case 'c':
					fmt.Print("\033[6n")
				default:
					move(1, 1)
					fmt.Printf("?? %x ", r)
				}
			case esc:
				if r == '[' {
					fmt.Print("CSI ")
					state = csi
				} else if r == 'O' {
					fmt.Print("SS3 ")
					state = ss3
				} else {
					fmt.Print("Final byte: ", r)
				}
			case ss3:
				switch r {
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
				}
				state = c0
			case csi:
				if r >= 0x30 && r <= 0x3f {
					n := 0
					for r >= 0x30 && r <= 0x39 {
						n = n*10 + int(r-0x30)
						scanner.Scan()
						r = scanner.Bytes()[0]
					}
					fmt.Print("Parameter: ", n)
					parms = append(parms, n)

					if r == ';' {
						fmt.Print("; ")
					}
				}
				if r >= 0x40 && r <= 0x7e {
					if r == 0x4d {
						fmt.Print("mouse ")
						state = mouse
					} else {
						state = c0
						switch r {
						case 0x41: // A
							fmt.Print("up")
						case 0x42: // B
							fmt.Print("down")
						case 0x43: // C
							fmt.Print("right")
						case 0x44: // D
							fmt.Print("left")
						case 0x46: // F
							fmt.Print("end")
						case 0x48: // H
							fmt.Print("home")
						case 0x4a: // J
							fmt.Print("clear")
						case 0x4b: // K
							fmt.Print("clear line")
						case 0x7e: // ~
							fmt.Print("special char ")
							switch parms[0] {
							case 1:
								fmt.Print("home")
							case 3:
								fmt.Print("del")
							case 4:
								fmt.Print("end")
							case 5:
								fmt.Print("page up")
							case 6:
								fmt.Print("page down")
							case 11:
								fmt.Print("F1")
							case 12:
								fmt.Print("F2")
							case 13:
								fmt.Print("F3")
							case 14:
								fmt.Print("F4")
							case 15:
								fmt.Print("F5")
							case 17:
								fmt.Print("F6")
							case 18:
								fmt.Print("F7")
							case 19:
								fmt.Print("F8")
							case 20:
								fmt.Print("F9")
							case 21:
								fmt.Print("F10")
							case 23:
								fmt.Print("F11")
							case 24:
								fmt.Print("F12")
							default:
								fmt.Print("special char ", parms)
							}
						default:
							fmt.Printf("Final byte: %c", r)
							state = c0
						}
					}
				}
			case mouse:
				switch r {
				case 0x20:
					fmt.Printf("left,")
				case 0x21:
					fmt.Printf("middle,")
				case 0x22:
					fmt.Printf("right,")
				case 0x23:
					fmt.Printf("up,")
				case 0x43:
					fmt.Printf("drag,")
				case 0x60:
					fmt.Printf("scroll up,")
				case 0x61:
					fmt.Printf("scroll down,")
				case 0x62:
					fmt.Printf("scroll left,")
				case 0x63:
					fmt.Printf("scroll right,")
				default:
					fmt.Printf("%x", r)
				}
				scanner.Scan()
				r5 := scanner.Bytes()
				fmt.Printf("x=%d,", r5[0]-32)

				scanner.Scan()
				r6 := scanner.Bytes()
				fmt.Printf("y=%d", r6[0]-32)

				state = c0
			}
		}
	}
}
