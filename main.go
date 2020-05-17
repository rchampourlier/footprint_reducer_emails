//+build !test

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"footprint_reducer_emails/controller"
	"footprint_reducer_emails/ui"
)

func main() {
	i, err := ui.NewGocuiUI()
	if err != nil {
		log.Panicln(err)
	}
	defer i.Close()

	server := os.Getenv("SERVER")
	username := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	c := controller.NewControllerWithCredentials(i, server, username, password)
	go func() {
		err := c.Run()
		if err != nil {
			panic(err)
		}
	}()

	i.Start()
}
