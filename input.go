package main

import (
	"bufio"
	"fmt"
	"os"
)

func input() WindowMessage {
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
				default:
					return WindowMessage{Action: "KeyPress", Key: string(r)}
				}
			case esc:
				if r == '[' {
					fmt.Print("CSI ")
					state = csi
				} else if r == 'O' {
					fmt.Print("SS3 ")
					state = ss3
				} else {
					move(1, 1)
					fmt.Print("Final byte: ", r)
				}
			case ss3:
				switch r {
				case 0x48:
					return WindowMessage{Action: "KeyPress", Key: "home"}
				case 0x46:
					return WindowMessage{Action: "KeyPress", Key: "end"}
				case 0x50:
					return WindowMessage{Action: "KeyPress", Key: "F1"}
				case 0x51:
					return WindowMessage{Action: "KeyPress", Key: "F2"}
				case 0x52:
					return WindowMessage{Action: "KeyPress", Key: "F3"}
				case 0x53:
					return WindowMessage{Action: "KeyPress", Key: "F4"}
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
					parms = append(parms, n)
				}
				if r >= 0x40 && r <= 0x7e {
					if r == 0x4d {
						state = mouse
					} else {
						state = c0
						switch r {
						case 0x41: // A
							return WindowMessage{Action: "KeyPress", Key: "up"}
						case 0x42: // B
							return WindowMessage{Action: "KeyPress", Key: "down"}
						case 0x43: // C
							return WindowMessage{Action: "KeyPress", Key: "right"}
						case 0x44: // D
							return WindowMessage{Action: "KeyPress", Key: "left"}
						case 0x46: // F
							return WindowMessage{Action: "KeyPress", Key: "end"}
						case 0x48: // H
							return WindowMessage{Action: "KeyPress", Key: "home"}
						case 0x4a: // J
							return WindowMessage{Action: "ClearScreen"}
						case 0x4b: // K
							return WindowMessage{Action: "ClearLine"}
						case 0x7e: // ~
							switch parms[0] {
							case 1:
								return WindowMessage{Action: "KeyPress", Key: "home"}
							case 3:
								return WindowMessage{Action: "KeyPress", Key: "delete"}
							case 4:
								return WindowMessage{Action: "KeyPress", Key: "end"}
							case 5:
								return WindowMessage{Action: "KeyPress", Key: "page up"}
							case 6:
								return WindowMessage{Action: "KeyPress", Key: "page down"}
							case 11:
								return WindowMessage{Action: "KeyPress", Key: "F1"}
							case 12:
								return WindowMessage{Action: "KeyPress", Key: "F2"}
							case 13:
								return WindowMessage{Action: "KeyPress", Key: "F3"}
							case 14:
								return WindowMessage{Action: "KeyPress", Key: "F4"}
							case 15:
								return WindowMessage{Action: "KeyPress", Key: "F5"}
							case 17:
								return WindowMessage{Action: "KeyPress", Key: "F6"}
							case 18:
								return WindowMessage{Action: "KeyPress", Key: "F7"}
							case 19:
								return WindowMessage{Action: "KeyPress", Key: "F8"}
							case 20:
								return WindowMessage{Action: "KeyPress", Key: "F9"}
							case 21:
								return WindowMessage{Action: "KeyPress", Key: "F10"}
							case 23:
								return WindowMessage{Action: "KeyPress", Key: "F11"}
							case 24:
								return WindowMessage{Action: "KeyPress", Key: "F12"}
							default:
								return WindowMessage{Action: "KeyPress", Key: fmt.Sprintf("special char %d", parms[0])}
							}
						default:
							fmt.Printf("Final byte: %c", r)
							state = c0
						}
					}
				}
			case mouse:
				scanner.Scan()
				x := scanner.Bytes()[0] - 32

				scanner.Scan()
				y := scanner.Bytes()[0] - 32

				switch r {
				case 0x20:
					return WindowMessage{Action: "MouseLeftButtonDown", X: int(x), Y: int(y)}
				case 0x21:
					return WindowMessage{Action: "MouseMiddleButtonDown", X: int(x), Y: int(y)}
				case 0x22:
					return WindowMessage{Action: "MouseRightButtonDown", X: int(x), Y: int(y)}
				case 0x23:
					return WindowMessage{Action: "MouseRelease", X: int(x), Y: int(y)}
				case 0x43:
					return WindowMessage{Action: "MouseDrag", X: int(x), Y: int(y)}
				case 0x60:
					return WindowMessage{Action: "MouseScrollUp", X: int(x), Y: int(y)}
				case 0x61:
					return WindowMessage{Action: "MouseScrollDown", X: int(x), Y: int(y)}
				case 0x62:
					return WindowMessage{Action: "MouseScrollLeft", X: int(x), Y: int(y)}
				case 0x63:
					return WindowMessage{Action: "MouseScrollRight", X: int(x), Y: int(y)}
				default:
					fmt.Printf("%x", r)
				}

				state = c0
			}
		}
	}
}
