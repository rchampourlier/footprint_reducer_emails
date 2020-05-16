//+build !test

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

// getServer
func getServer(g *gocui.Gui, ch chan string) {
	g.SetManagerFunc(func(gui *gocui.Gui) error {
		if err := setGlobalKeybindings(gui); err != nil {
			return err
		}
		if err := getUserInput(gui, "Enter the server domain and port (e.g. imap.gmail.com:993):", ch); err != nil {
			return err
		}
		return nil
	})
}

// getUsername
func getUsername(g *gocui.Gui, ch chan string) {
	g.SetManagerFunc(func(gui *gocui.Gui) error {
		if err := setGlobalKeybindings(gui); err != nil {
			return err
		}
		if err := getUserInput(gui, "Enter your IMAP username (generally your email address):", ch); err != nil {
			return err
		}
		return nil
	})
}

// getServer
func getPassword(g *gocui.Gui, ch chan string) {
	g.SetManagerFunc(func(gui *gocui.Gui) error {
		if err := setGlobalKeybindings(gui); err != nil {
			return err
		}
		if err := getUserInput(gui, "Enter your IMAP password:", ch); err != nil {
			return err
		}
		return nil
	})
}

// getUserInput displays a view to fetch the user's input.
func getUserInput(g *gocui.Gui, msg string, ch chan string) error {
	g.SetManagerFunc(func(gui *gocui.Gui) error {
		maxX, maxY := g.Size()
		inputView, err := g.SetView("input", maxX/2-30, maxY/2+2, maxX/2+30, maxY/2+4)
		if err != nil {
			if err != gocui.ErrUnknownView {
				log.Panicf("error setting input view: %x", err)
			}
			inputView.Title = msg
			inputView.Editable = true
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

func displayInformation(g *gocui.Gui, server, username, password string) {
	g.SetManagerFunc(func(gui *gocui.Gui) error {
		setGlobalKeybindings(gui)
		maxX, maxY := g.Size()
		v, err := g.SetView("information", maxX/2-30, maxY/2+1, maxX/2+30, maxY/2+5)
		v.Title = "Information:"
		if err != nil {
			if err != gocui.ErrUnknownView {
				return fmt.Errorf("error setting information view: %w", err)
			}
			if _, err := g.SetCurrentView("information"); err != nil {
				return fmt.Errorf("error setting current view to information: %w", err)
			}
			fmt.Fprintf(v, "Server: %s\nUsername: %s\nPassword: %s\n", server, username, password)
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

type program struct {
	gui       *gocui.Gui
	serverURL string
	username  string
	password  string
}

func newProgram(g *gocui.Gui) *program {
	return &program{g, "", "", ""}
}

func (p *program) run() error {
	ch := make(chan string, 0)

	// Get server URL
	getServer(p.gui, ch)
	server := <-ch

	// Get email username
	getUsername(p.gui, ch)
	username := <-ch

	// Get email password
	getPassword(p.gui, ch)
	password := <-ch

	// Closing the channel
	close(ch)

	// Display information
	displayInformation(p.gui, server, username, password)

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true

	p := newProgram(g)
	go func() {
		err := p.run()
		if err != nil {
			// handle error
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
