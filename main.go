//+build !test

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"footprint_reducer_emails/controller"
	"footprint_reducer_emails/emailclient"
	"footprint_reducer_emails/ui"
)

func main() {
	i, err := ui.NewGocuiUI()
	if err != nil {
		log.Panicln(err)
	}
	defer i.Close()

	w := emailclient.NewMockClientWrapper()
	//w := emailclient.NewImapClientWrapper()
	c := controller.NewController(w, i)
	go func() {
		err := c.Run()
		if err != nil {
			panic(err)
		}
	}()

	i.Start() // the interface has the control
}
