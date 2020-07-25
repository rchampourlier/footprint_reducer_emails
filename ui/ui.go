package ui

import (
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

// GocuiUI is an UI implementation using gocui
type GocuiUI struct {
	gui *gocui.Gui
}

// UI is the interface for an user interface
type UI interface {
	Start()
	Close()

	// Views
	Information(title, message string)
	List(items []string, ch chan<- Event)
	Page(title, content string)
	StringWithMaskInput(msg string, mask rune, ch chan<- Event) error
	StringInput(msg string, ch chan<- Event) error
}

// EventType is an enum representing the type of UI Event
// of the struct.
//
// EventTypeStringInputReturned: event generated when an input view
//   returns a value. The value is a string contained in the
//   `Data` field.
type EventType int

const (
	// EventTypeStringInputReturned indicates a value has been returned
	EventTypeStringInputReturned EventType = iota

	// EventTypeItemSelected indicates an item was selected in the view
	EventTypeItemSelected EventType = iota

	// EventError indicates an error in the `Event`
	EventError = iota
)

// Event represents an event generated by the user interface
type Event struct {
	Type EventType
	Data interface{}
	Err  error
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

// Close closes the interface if needed
func (g *GocuiUI) Close() {
	g.gui.Close()
}

// Views
// =====

// Information displays an information box with the specified title and message.
func (g *GocuiUI) Information(title, message string) {
	g.gui.SetManagerFunc(func(gui *gocui.Gui) error {
		maxX, maxY := g.gui.Size()
		v, err := g.gui.SetView("information", maxX/2-30, maxY/2+1, maxX/2+30, maxY/2+5)
		v.Title = title
		if err != nil {
			if err != gocui.ErrUnknownView {
				return fmt.Errorf("error setting information view: %w", err)
			}
			if _, err := g.gui.SetCurrentView("information"); err != nil {
				return fmt.Errorf("error setting current view to information: %w", err)
			}
			fmt.Fprintf(v, message)
		}
		return nil
	})

	g.setGlobalKeybindings()
}

// List displays a list view where each line is an item of the passed
// `items` slice.
func (g *GocuiUI) List(items []string, evtCh chan<- Event) {
	viewChan := make(chan *gocui.View, 0) // to send the view to the goroutine

	g.gui.SetManagerFunc(func(gui *gocui.Gui) error {
		maxX, maxY := g.gui.Size()
		v, err := g.gui.SetView("list", 0, 0, maxX-1, maxY-1)
		if err != nil {
			if err == gocui.ErrUnknownView {
				v.Frame = false
				v.Highlight = true
				v.SelBgColor = gocui.ColorGreen
				v.SelFgColor = gocui.ColorBlack

				// Sending the viewChan to the goroutine displaying content
				viewChan <- v
			} else {
				return fmt.Errorf("error creating `List` view: %w", err)
			}
		}
		if _, err := g.gui.SetCurrentView("list"); err != nil {
			return fmt.Errorf("error setting `list` view to current: %w", err)
		}

		return nil
	})

	// This goroutine displays the content of the list.
	go func() {
		v := <-viewChan
		for _, item := range items {
			fmt.Fprintf(v, "%s\n", item)
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
	g.gui.SetKeybinding("list", gocui.KeyArrowDown, gocui.ModNone, scrollDownOneLine)
	g.gui.SetKeybinding("list", gocui.MouseWheelDown, gocui.ModNone, scrollDownOneLine)

	// Page down
	scrollDownOnePage := func(_ *gocui.Gui, v *gocui.View) error {
		x, y := v.Origin()
		_, rows := v.Size()
		v.SetOrigin(x, y+rows)
		return nil
	}
	g.gui.SetKeybinding("list", gocui.KeyArrowRight, gocui.ModNone, scrollDownOnePage)

	// Line up
	scrollUpOneLine := func(_ *gocui.Gui, v *gocui.View) error {
		x, y := v.Origin()
		v.SetOrigin(x, y-1)
		return nil
	}
	g.gui.SetKeybinding("list", gocui.KeyArrowUp, gocui.ModNone, scrollUpOneLine)
	g.gui.SetKeybinding("list", gocui.MouseWheelUp, gocui.ModNone, scrollUpOneLine)

	// Page up
	scrollUpOnePage := func(_ *gocui.Gui, v *gocui.View) error {
		x, y := v.Origin()
		_, rows := v.Size()
		v.SetOrigin(x, y-rows)
		return nil
	}
	g.gui.SetKeybinding("list", gocui.KeyArrowLeft, gocui.ModNone, scrollUpOnePage)

	// Return an `ItemSelected` event on <ENTER>
	g.gui.SetKeybinding("list", gocui.KeyEnter, gocui.ModNone, func(_ *gocui.Gui, v *gocui.View) error {
		// read the highlighted line
		// return a command to display the messages for the highlighted sender
		_, lineIndex := v.Origin()
		evtCh <- Event{EventTypeItemSelected, lineIndex, nil}
		return nil
	})

	g.setGlobalKeybindings()
}

// Page displays the specified content on the full space of the screen,
// below the specified title.
func (g *GocuiUI) Page(title, content string) {
	g.gui.SetManagerFunc(func(gui *gocui.Gui) error {
		maxX, maxY := g.gui.Size()
		v, err := g.gui.SetView("page", 0, 0, maxX, maxY)
		v.Title = title
		v.Editable = false
		if err != nil {
			if err != gocui.ErrUnknownView {
				return fmt.Errorf("error setting information view: %w", err)
			}
			if _, err := g.gui.SetCurrentView("page"); err != nil {
				return fmt.Errorf("error setting current view to `page`: %w", err)
			}
			fmt.Fprintf(v, content)
		}
		return nil
	})

	g.setGlobalKeybindings()
}

func (g *GocuiUI) setGlobalKeybindings() error {
	if err := g.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {

		return err
	}
	return nil
}

func (g *GocuiUI) StringInput(msg string, ch chan<- Event) error {
	return g._getUserInput(msg, false, ' ', ch)
}

func (g *GocuiUI) StringWithMaskInput(msg string, mask rune, ch chan<- Event) error {
	return g._getUserInput(msg, true, mask, ch)
}

// _getUserInput is the internal function used to display a view to fetch the user's input
// with or without a mask.
//   - `withMask`: set the view with `mask` as mask. If false, the `mask` value is ignored
//   - `mask`: the mask to use for the view if `withMask` is true
func (g *GocuiUI) _getUserInput(msg string, withMask bool, mask rune, ch chan<- Event) error {
	g.gui.SetManagerFunc(func(gui *gocui.Gui) error {
		maxX, maxY := g.gui.Size()
		inputView, err := g.gui.SetView("input", maxX/2-30, maxY/2+2, maxX/2+30, maxY/2+4)
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
			if _, err := g.gui.SetCurrentView("input"); err != nil {
				return fmt.Errorf("error setting `input` view current: %w", err)
			}
		}
		return nil
	})
	if err := g.setGlobalKeybindings(); err != nil {
		return err
	}
	err := g.gui.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			line, err := v.Line(0)
			if err != nil {
				if err.Error() == "invalid point" {
					ch <- Event{EventTypeStringInputReturned, "", err}
					return nil
				}
				ch <- Event{EventTypeStringInputReturned, "", err}
				return err
			}
			ch <- Event{EventTypeStringInputReturned, line, nil}
			return nil
		},
	)
	return err
}

func logger() *log.Logger {
	f, err := os.Create("./log.txt")
	if err != nil {
		panic(err)
	}
	logger := log.New(f, "", 0)
	return logger
}
