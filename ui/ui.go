package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// GocuiUI is an UI implementation using gocui
type GocuiUI struct {
	gui *gocui.Gui
}

// UI is the interface for an user interface
type UI interface {
	GetServer(ch chan<- string)
	GetUsername(ch chan<- string)
	GetPassword(ch chan<- string)
	DisplayInformation(server, username, password string)
	MainView() chan<- string
	Start()
	Close()
}

// NewGocuiUI creates a new gocui user interface
func NewGocuiUI() (*GocuiUI, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	g.Cursor = true
	if err != nil {
		return nil, err
	}
	return &GocuiUI{g}, nil
}

// Start starts the interface
//
// The function is blocking and only stops on interface error or
// when an interruption signal is received.
func (g *GocuiUI) Start() {
	if err := g.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

// GetServer displays an interface to retrieve the email
// server.
func (g *GocuiUI) GetServer(ch chan<- string) {
	g.gui.SetManagerFunc(func(gui *gocui.Gui) error {
		if err := setGlobalKeybindings(gui); err != nil {
			return err
		}
		if err := getUserInput(gui, "Enter the server domain and port (e.g. imap.gmail.com:993):", ch); err != nil {
			return err
		}
		return nil
	})
}

// GetUsername dispays an interface to fetch the IMAP username
func (g *GocuiUI) GetUsername(ch chan<- string) {
	g.gui.SetManagerFunc(func(gui *gocui.Gui) error {
		if err := setGlobalKeybindings(gui); err != nil {
			return err
		}
		if err := getUserInput(gui, "Enter your IMAP username (generally your email address):", ch); err != nil {
			return err
		}
		return nil
	})
}

// GetPassword dispays an interface to fetch the IMAP password
func (g *GocuiUI) GetPassword(ch chan<- string) {
	g.gui.SetManagerFunc(func(gui *gocui.Gui) error {
		if err := getUserInputWithMask(gui, "Enter your IMAP password:", '•', ch); err != nil {
			return err
		}
		return nil
	})
	setGlobalKeybindings(g.gui)
}

// MainView displays the main interface with the list of senders
func (g *GocuiUI) MainView() chan<- string {
	linesChan := make(chan string, 0)     // to send lines to display in the view
	viewChan := make(chan *gocui.View, 0) // to send the view to the goroutine

	// send v to the goroutine using a chan
	g.gui.SetManagerFunc(func(gui *gocui.Gui) error {
		maxX, maxY := g.gui.Size()
		v, err := g.gui.SetView("main", 0, 0, maxX-1, maxY-1)
		if err != nil {
			if err == gocui.ErrUnknownView {
				v.Frame = false
				v.Highlight = true
				v.SelBgColor = gocui.ColorGreen
				v.SelFgColor = gocui.ColorBlack

				// Sending the viewChan to the goroutine displaying content
				viewChan <- v
			} else {
				return fmt.Errorf("error creating \"MainView\": %w", err)
			}
		}
		if _, err := g.gui.SetCurrentView("main"); err != nil {
			return fmt.Errorf("error setting main view to current: %w", err)
		}
		return nil
	})

	// This goroutine will display the content of the main view.
	// Each line passed over `viewChan` gets added to the view's buffer.
	go func() {
		v := <-viewChan
		for line := range linesChan {
			fmt.Fprintf(v, "%s\n", line)
			g.gui.Update(func(g *gocui.Gui) error {
				// Empty update to trigger a display of the updated
				// view's buffer.
				// The buffer update (`fmt.Fprintf...` above) must
				// not be in there, otherwise it may get called
				// several times.
				return nil
			})
		}
	}()

	// Line down
	scrollDownOneLine := func(_ *gocui.Gui, v *gocui.View) error {
		x, y := v.Origin()
		v.SetOrigin(x, y+1)
		return nil
	}
	g.gui.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, scrollDownOneLine)
	g.gui.SetKeybinding("main", gocui.MouseWheelDown, gocui.ModNone, scrollDownOneLine)

	// Page down
	scrollDownOnePage := func(_ *gocui.Gui, v *gocui.View) error {
		x, y := v.Origin()
		_, rows := v.Size()
		v.SetOrigin(x, y+rows)
		return nil
	}
	g.gui.SetKeybinding("main", gocui.KeyArrowRight, gocui.ModNone, scrollDownOnePage)

	// Line up
	scrollUpOneLine := func(_ *gocui.Gui, v *gocui.View) error {
		x, y := v.Origin()
		v.SetOrigin(x, y-1)
		return nil
	}
	g.gui.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, scrollUpOneLine)
	g.gui.SetKeybinding("main", gocui.MouseWheelUp, gocui.ModNone, scrollUpOneLine)

	// Page up
	scrollUpOnePage := func(_ *gocui.Gui, v *gocui.View) error {
		x, y := v.Origin()
		_, rows := v.Size()
		v.SetOrigin(x, y-rows)
		return nil
	}
	g.gui.SetKeybinding("main", gocui.KeyArrowLeft, gocui.ModNone, scrollUpOnePage)

	setGlobalKeybindings(g.gui)
	return linesChan
}

// Close closes the interface if needed
func (g *GocuiUI) Close() {
	g.gui.Close()
}

// DisplayInformation displays the server information
func (g *GocuiUI) DisplayInformation(server, username, password string) {
	g.gui.SetManagerFunc(func(gui *gocui.Gui) error {
		maxX, maxY := g.gui.Size()
		v, err := g.gui.SetView("information", maxX/2-30, maxY/2+1, maxX/2+30, maxY/2+5)
		v.Title = "Information:"
		if err != nil {
			if err != gocui.ErrUnknownView {
				return fmt.Errorf("error setting information view: %w", err)
			}
			if _, err := g.gui.SetCurrentView("information"); err != nil {
				return fmt.Errorf("error setting current view to information: %w", err)
			}
			fmt.Fprintf(v, "Server: %s\nUsername: %s\nPassword: %s\n", server, username, "••••••••")
		}
		return nil
	})
	setGlobalKeybindings(g.gui)
}

func setGlobalKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func getUserInput(g *gocui.Gui, msg string, ch chan<- string) error {
	return _getUserInput(g, msg, false, ' ', ch)
}

func getUserInputWithMask(g *gocui.Gui, msg string, mask rune, ch chan<- string) error {
	return _getUserInput(g, msg, true, mask, ch)
}

// _getUserInput is the internal function used to display a view to fetch the user's input
// with or without a mask.
//   - `withMask`: set the view with `mask` as mask. If false, the `mask` value is ignored
//   - `mask`: the mask to use for the view if `withMask` is true
func _getUserInput(g *gocui.Gui, msg string, withMask bool, mask rune, ch chan<- string) error {
	g.SetManagerFunc(func(gui *gocui.Gui) error {
		maxX, maxY := g.Size()
		inputView, err := g.SetView("input", maxX/2-30, maxY/2+2, maxX/2+30, maxY/2+4)
		if err != nil {
			if err == gocui.ErrUnknownView {
				inputView.Title = msg
				inputView.Editable = true
				if withMask {
					inputView.Mask = mask
				}
			} else {
				return fmt.Errorf("error creating `input` view: %w", err)
			}
			if _, err := g.SetCurrentView("input"); err != nil {
				return fmt.Errorf("error setting `input` view current: %w", err)
			}
		}
		return nil
	})
	if err := setGlobalKeybindings(g); err != nil {
		return err
	}
	err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			// TODO get called twice, panics the 2nd time
			line, err := v.Line(0)
			if err != nil {
				if err.Error() == "invalid point" {
					ch <- ""
					return nil
				}
				ch <- ""
				return err
			}
			ch <- line
			return nil
		},
	)
	return err
}
