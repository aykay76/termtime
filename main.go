package main

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
	wm.AddWindow(win2)
	wm.AddWindow(win)
	wm.Start()
}
