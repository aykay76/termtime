package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

// WindowMessage is a message that can be sent from the window manager to a listening channel
type WindowMessage struct {
	Window *Window
	Action string
	Key    string
	X      int
	Y      int
}

// Window represents a window in the application
type Window struct {
	X, Y, Width, Height int
	Content             []string
	// add hierarchy for stacked windows
	Children []*Window
	Parent   *Window
	// add a flag to indicate whether the window is clickable, we don't want the root window to be brought to the front
	Clickable bool
	// add a flag to indicate whether the border is visible
	Border bool
}

func NewWindow(x, y, width, height int, border bool, content []string) *Window {
	return &Window{
		X:       x,
		Y:       y,
		Width:   width,
		Height:  height,
		Content: content,
		Border:  border,
	}
}

type WindowManager struct {
	ScreenWidth, ScreenHeight int

	// Windows is a slice that holds pointers to Window objects, representing multiple windows in the application.
	Windows       []*Window
	FocusedWindow *Window

	ForceRedraw bool

	// store old state of terminal
	oldState *term.State
}

func NewWindowManager() *WindowManager {
	screenWidth, screenHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println(err)
	}

	wm := &WindowManager{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		Windows:      []*Window{},
	}

	// we start with a single window that covers the entire screen
	// if we want to add a window on top of this, we can add it to the window manager
	// and it will be rendered on top of the root window
	screen := NewWindow(1, 1, screenWidth-1, screenHeight-1, false, []string{""})
	wm.Windows = append(wm.Windows, screen)

	// setup the terminal
	wm.oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
	}
	smcup()
	clear()
	hideCursor()
	echoOff()
	enableMouseTracking()

	return wm
}

// add window to window manager, and add window to root window or parent window
func (wm *WindowManager) AddWindow(win *Window) {
	// add window to window manager stack
	wm.Windows = append(wm.Windows, win)
	// set focus on the new window
	wm.FocusedWindow = win
	// render the window manager
	wm.Render()
}

// TODO: add mechanism to detect dirty cells and only render those cells
// TODO: add support for colors in content

func (wm *WindowManager) renderWindow(window *Window, screen [][]rune, mask [][]int) {
	// Mark the areas covered by children in the mask
	for _, child := range window.Children {
		for i := 0; i < child.Height; i++ {
			if child.Y+i >= window.Height {
				break
			}
			for j := 0; j < child.Width; j++ {
				if child.X+j >= window.Width {
					break
				}
				mask[child.Y+i][child.X+j]++
			}
		}
	}

	// Mark the areas covered by other windows in the stack but not this window
	// TODO: this is a naive implementation - ignore windows before this window
	// ignore windows "below" this window and ignore this window itself
	found := false
	for _, win := range wm.Windows {
		// ignore this window
		if win == window {
			found = true
			continue
		}

		// ignore windows "below" this window
		if !found {
			continue
		}

		// mark the areas covered by this window
		for i := 0; i < win.Height; i++ {
			if win.Y+i >= wm.ScreenHeight {
				break
			}
			for j := 0; j < win.Width; j++ {
				if win.X+j >= wm.ScreenWidth {
					break
				}
				if win.Y+i >= window.Y && win.Y+i < window.Y+window.Height &&
					win.X+j >= window.X && win.X+j < window.X+window.Width {
					mask[win.Y+i][win.X+j]++
				}
			}
		}
	}

	// Render the window content where it's not obscured by children
	for i := 0; i < window.Height; i++ {
		if window.Y+i >= wm.ScreenHeight {
			break
		}
		for j := 0; j < window.Width; j++ {
			if window.X+j >= wm.ScreenWidth {
				break
			}

			// draw a border around the window
			if window.Border {
				if i == 0 {
					if j == 0 {
						screen[window.Y+i][window.X+j] = '┌'
					} else if j == window.Width-1 {
						screen[window.Y+i][window.X+j] = '┐'
					} else {
						screen[window.Y+i][window.X+j] = '─'
					}
				} else if i == window.Height-1 {
					if j == 0 {
						screen[window.Y+i][window.X+j] = '└'
					} else if j == window.Width-1 {
						screen[window.Y+i][window.X+j] = '┘'
					} else {
						screen[window.Y+i][window.X+j] = '─'
					}
				} else if j == 0 || j == window.Width-1 {
					screen[window.Y+i][window.X+j] = '│'
				}
			}

			// at this point I need to move the cursor ready to output the content
			move(window.X+j, window.Y+i)

			// also need to clip the window content to the height and width of the window
			if mask[i][j] == 0 {
				contentOffset := 0
				if window.Border {
					contentOffset = 1
				}
				if i-contentOffset >= 0 && i-contentOffset < len(window.Content) && j-contentOffset >= 0 && j-contentOffset < len(window.Content[i-contentOffset]) {
					if window.Border && (i == 0 || i == window.Height-1 || j == 0 || j == window.Width-1) {
						continue
					}
					printCenter(1, window.Content[i-contentOffset])
					screen[window.Y+i][window.X+j] = rune(window.Content[i-contentOffset][j-contentOffset])
				} else {
					if window.Border && (i == 0 || i == window.Height-1 || j == 0 || j == window.Width-1) {
						continue
					}
					screen[window.Y+i][window.X+j] = ' '
				}
			}
		}
	}

	// Render the children windows
	for _, child := range window.Children {
		wm.renderWindow(child, screen, mask)
	}

	// render the screen
	for i := range screen {
		move(1, i)
		for j := range screen[i] {
			// setBackground(mask[i][j])
			fmt.Print(string(screen[i][j]))
		}
	}

	move(1, 1)
}

func (wm *WindowManager) Start(c chan WindowMessage) {
	wm.Render()

	// handle keyboard and mouse input
	go func() {
		for {
			message := input()

			// find the window that this message occurred on, for mouse events it will be wherever the cursor is
			// for keyboard events it will be the window that has focus
			if message.X > 0 && message.Y > 0 {
				for i := len(wm.Windows) - 1; i >= 0; i-- {
					window := wm.Windows[i]
					if message.X >= window.X && message.X < window.X+window.Width &&
						message.Y >= window.Y && message.Y < window.Y+window.Height {
						message.Window = window

						printCenterf(3, "window: %#v", window)
						if message.Action == "MouseLeftButtonDown" {
							wm.Windows = append(wm.Windows[:i], wm.Windows[i+1:]...)
							wm.Windows = append(wm.Windows, message.Window)
							wm.FocusedWindow = message.Window
							wm.ForceRedraw = true
						}

						break
					}
				}
			} else {
				// find the window that has focus
				message.Window = wm.FocusedWindow
			}

			printCenterf(2, "message: %#v", message)
			c <- message
		}
	}()

	// loop forever checking for keyboard and mouse input
	for {
		// redraw the screen if the terminal size has changed
		time.Sleep(10 * time.Millisecond)

		// get the terminal size
		width, height, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			fmt.Println(err)
		}

		// if the terminal size has changed, redraw the screen
		if wm.ForceRedraw || width != wm.ScreenWidth || height != wm.ScreenHeight {
			printCenterf(wm.ScreenHeight-1, "width: %d, height: %d\n", width, height)
			wm.ScreenWidth = width
			wm.ScreenHeight = height
			wm.Render()
			wm.ForceRedraw = false
		}
	}
}

func (wm *WindowManager) Render() {
	screen := make([][]rune, wm.ScreenHeight)
	for i := range screen {
		screen[i] = make([]rune, wm.ScreenWidth)
		for j := range screen[i] {
			screen[i][j] = ' '
		}
	}

	mask := make([][]int, wm.ScreenHeight)
	for i := range mask {
		mask[i] = make([]int, wm.ScreenWidth)
	}

	// Render the windows
	for _, window := range wm.Windows {
		wm.renderWindow(window, screen, mask)
	}
}

func (wm *WindowManager) Close() {
	term.Restore(int(os.Stdin.Fd()), wm.oldState)
	showCursor()
	echoOn()
	disableMouseTracking()
	rmcup()
}
