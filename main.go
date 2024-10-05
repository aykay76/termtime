package main

import (
	"os"
)

var wm *WindowManager

func main() {
	// create a window manager and add some example windows
	wm = NewWindowManager()
	win := NewWindow(10, 10, 20, 10, true, []string{
		"This is the first window",
		"It has multiple lines",
		"And is positioned at 10, 10",
		"It is 20 characters wide and 10 characters tall",
		"It is the first window added to the window manager",
		"Press 'q' to quit",
		"Press 'w' to add a new window",
		"Press 'e' to remove the last window",
		"Press 'r' to remove all windows",
		"Press 't' to toggle window visibility",
		"Press 'y' to toggle window focus",
	})
	win2 := NewWindow(15, 15, 20, 10, true, []string{
		"This window should sit on top of the first window",
		"It is positioned at 15, 15",
		"It is 20 characters wide and 10 characters tall",
		"It is the second window added to the window manager",
		"It has a border",
		"It is on top of the stack",
	})
	wm.AddWindow(win)
	wm.AddWindow(win2)

	// create a channel that i can receive messages from the window manager on
	c := make(chan WindowMessage, 1)

	go func() {
		for {
			msg := <-c
			switch msg.Action {
			case "MouseMove":
				// do something with the mouse move event
				printCenterf(1, "Mouse moved to %d,%d", msg.X, msg.Y)
			case "MouseLeftButtonDown":
				// do something with the mouse click event
			case "KeyPress":
				// do something with the key press event
				switch msg.Key {
				case "w":
					win3 := NewWindow(18, 18, 20, 10, true, []string{
						"This is the third window",
						"It has multiple lines",
						"And is positioned at 18, 18",
						"It is 20 characters wide and 10 characters tall",
						"It is the third window added to the window manager",
						"It is on top of the stack",
					})
					wm.AddWindow(win3)
				case "q":
					wm.Close()
					os.Exit(0)
				}
			}
		}
	}()

	wm.Start(c)
}
