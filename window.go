package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type Window struct {
	X, Y, Width, Height int
	Content             []string
	// add hierarchy for stacked windows
	Children []*Window
	Parent   *Window
}

func NewWindow(x, y, width, height int, content []string) *Window {
	return &Window{
		X:       x,
		Y:       y,
		Width:   width,
		Height:  height,
		Content: content,
	}
}

type WindowManager struct {
	ScreenWidth, ScreenHeight int
	RootWindow                *Window
}

func NewWindowManager() *WindowManager {
	screenWidth, screenHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println(err)
	}

	wm := &WindowManager{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		RootWindow:   NewWindow(0, 0, screenWidth, screenHeight, []string{}),
	}

	return wm
}

// add window to window manager, and add window to root window or parent window
func (wm *WindowManager) AddWindow(win *Window, parentWindow *Window) {
	if parentWindow != nil {
		parentWindow.Children = append(parentWindow.Children, win)
		win.Parent = parentWindow
	} else {
		// add window to root window if no parentWindow is provided
		wm.RootWindow.Children = append(wm.RootWindow.Children, win)
		win.Parent = wm.RootWindow
	}
}

// TODO: only render parts of the window that have changed
// TODO: add support for colors

func (wm *WindowManager) renderWindow(window *Window, screen [][]rune) {
	// Create a mask to track which parts of the window are obscured by children
	mask := make([][]bool, window.Height)
	for i := range mask {
		mask[i] = make([]bool, window.Width)
	}

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
				mask[child.Y+i][child.X+j] = true
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
			if !mask[i][j] && i < len(window.Content) && j < len(window.Content[i]) {
				screen[window.Y+i][window.X+j] = rune(window.Content[i][j])
			}
		}
	}

	// Render the children windows
	for _, child := range window.Children {
		wm.renderWindow(child, screen)
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

	wm.renderWindow(wm.RootWindow, screen)
}
