package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

// GocuiUI is an UI implementation using gocui
type GocuiUI struct {
	gui *gocui.Gui
}

// UI is the interface for an user interface
type UI interface {
	GetServer(ch chan string)
	GetUsername(ch chan string)
	GetPassword(ch chan string)
	DisplayInformation(server, username, password string)
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
		log.Panicln(err)
	}
}

// GetServer displays an interface to retrieve the email
// server.
func (g *GocuiUI) GetServer(ch chan string) {
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

// GetUsername
func (g *GocuiUI) GetUsername(ch chan string) {
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

// GetPassword
func (g *GocuiUI) GetPassword(ch chan string) {
	g.gui.SetManagerFunc(func(gui *gocui.Gui) error {
		if err := setGlobalKeybindings(gui); err != nil {
			return err
		}
		if err := getUserInputWithMask(gui, "Enter your IMAP password:", '•', ch); err != nil {
			return err
		}
		return nil
	})
}

// Close closes the interface if needed
func (g *GocuiUI) Close() {
	g.gui.Close()
}

// _getUserInput is the internal function used to display a view to fetch the user's input
// with or without a mask.
//   - `withMask`: set the view with `mask` as mask. If false, the `mask` value is ignored
//   - `mask`: the mask to use for the view if `withMask` is true
func _getUserInput(g *gocui.Gui, msg string, withMask bool, mask rune, ch chan string) error {
	g.SetManagerFunc(func(gui *gocui.Gui) error {
		maxX, maxY := g.Size()
		inputView, err := g.SetView("input", maxX/2-30, maxY/2+2, maxX/2+30, maxY/2+4)
		if err != nil {
			if err != gocui.ErrUnknownView {
				log.Panicf("error setting input view: %x", err)
			}
			inputView.Title = msg
			inputView.Editable = true
			if withMask {
				inputView.Mask = mask
			}
			if _, err := g.SetCurrentView("input"); err != nil {
				log.Panicf("error setting current view to input: %x", err)
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

func getUserInput(g *gocui.Gui, msg string, ch chan string) error {
	return _getUserInput(g, msg, false, ' ', ch)
}

func getUserInputWithMask(g *gocui.Gui, msg string, mask rune, ch chan string) error {
	return _getUserInput(g, msg, true, mask, ch)
}

func (g *GocuiUI) DisplayInformation(server, username, password string) {
	g.gui.SetManagerFunc(func(gui *gocui.Gui) error {
		setGlobalKeybindings(gui)
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
