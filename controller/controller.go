package controller

import (
	"footprint_reducer_emails/ui"
)

// Controller represents a controller and stored the reference
// to the UI and the state of the program execution.
type Controller struct {
	ui        ui.UI
	serverURL string
	username  string
	password  string
}

// NewController returns a new controller with the specified UI
func NewController(i ui.UI) *Controller {
	return &Controller{i, "", "", ""}
}

// Run executes the program.
func (c *Controller) Run() error {
	ch := make(chan string, 0)

	// Get server URL
	c.ui.GetServer(ch)
	server := <-ch

	// Get email username
	c.ui.GetUsername(ch)
	username := <-ch

	// Get email password
	c.ui.GetPassword(ch)
	password := <-ch

	// Closing the channel
	close(ch)

	// Display information
	c.ui.DisplayInformation(server, username, password)

	return nil
}
