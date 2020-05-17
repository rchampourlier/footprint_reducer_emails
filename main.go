//+build !test

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"footprint_reducer_emails/ui"
)

type program struct {
	ui        ui.UI
	serverURL string
	username  string
	password  string
}

func newProgram(i ui.UI) *program {
	return &program{i, "", "", ""}
}

func (p *program) run() error {
	ch := make(chan string, 0)

	// Get server URL
	p.ui.GetServer(ch)
	server := <-ch

	// Get email username
	p.ui.GetUsername(ch)
	username := <-ch

	// Get email password
	p.ui.GetPassword(ch)
	password := <-ch

	// Closing the channel
	close(ch)

	// Display information
	p.ui.DisplayInformation(server, username, password)

	return nil
}

func main() {
	i, err := ui.NewGocuiUI()
	if err != nil {
		log.Panicln(err)
	}
	defer i.Close()

	p := newProgram(i)
	go func() {
		err := p.run()
		if err != nil {
			// handle error
		}
	}()

	i.Start()
}
