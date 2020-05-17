package controller

import (
	"fmt"

	"footprint_reducer_emails/emailclient"
	"footprint_reducer_emails/emailtools"
	"footprint_reducer_emails/ui"
)

// Controller represents a controller and stored the reference
// to the UI and the state of the program execution.
type Controller struct {
	ui       ui.UI
	server   string
	username string
	password string
}

// NewController returns a new controller with the specified UI
func NewController(i ui.UI) *Controller {
	return NewControllerWithCredentials(i, "", "", "")
}

// NewControllerWithCredentials initializes a new controller with
// the specified UI and server  and credentials.
func NewControllerWithCredentials(i ui.UI, server, username, password string) *Controller {
	return &Controller{i, server, username, password}
}

// Run executes the program.
func (c *Controller) Run() error {
	ch := make(chan string, 0)

	if c.server == "" {
		// Get server
		c.ui.GetServer(ch)
		c.server = <-ch
	}

	if c.username == "" {
		// Get email username
		c.ui.GetUsername(ch)
		c.username = <-ch
	}

	if c.password == "" {
		// Get email password
		c.ui.GetPassword(ch)
		c.password = <-ch
	}

	// Closing the channel
	close(ch)

	// Display information
	//c.ui.DisplayInformation(server, username, password)

	// Display the main view
	mainViewCh := c.ui.MainView()
	err := c.FetchEmails(mainViewCh)

	return err
}

// FetchEmails fetches emails from the IMAP server using the specified
//  and credentials.
func (c *Controller) FetchEmails(ch chan<- string) error {
	client, err := emailclient.ConnectAndLogin(c.server, c.username, c.password)
	if err != nil {
		return fmt.Errorf("failed to connect to IMAP server: %w", err)
	}

	// Don't forget to logout
	defer client.Logout()

	// TODO remove this constant
	const mailboxName = "[Gmail]/Tous les messages"

	messages, err := client.FetchMessages(mailboxName)
	stats := emailtools.StatsOnSenders(messages)
	emailtools.SortSendersStatBySize(stats)

	for i, stat := range stats {
		str := fmt.Sprintf("%04d | %s | %d messages | %d MB | %s", i, stat.Sender.Address(), stat.MessagesCount, stat.TotalSize/1024^2, stat.LatestMessageDate)
		ch <- str
	}
	close(ch)

	return err
}
