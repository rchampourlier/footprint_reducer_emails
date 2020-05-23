package controller

import (
	"fmt"

	"footprint_reducer_emails/emailclient"
	"footprint_reducer_emails/emailtools"
	uii "footprint_reducer_emails/ui"
)

// Controller represents a controller and stored the reference
// to the UI and the state of the program execution.
type Controller struct {
	ui       uii.UI
	server   string
	username string
	password string
}

// NewController returns a new controller with the specified UI
func NewController(i uii.UI) *Controller {
	return NewControllerWithCredentials(i, "", "", "")
}

// NewControllerWithCredentials initializes a new controller with
// the specified UI and server  and credentials.
func NewControllerWithCredentials(i uii.UI, server, username, password string) *Controller {
	return &Controller{i, server, username, password}
}

// Run executes the program.
func (c *Controller) Run() error {
	ui := c.ui
	uiEventCh := make(chan uii.Event, 0)
	defer close(uiEventCh)

	handleInput := func(inputFunc func(ch chan<- uii.Event)) (string, error) {
		inputFunc(uiEventCh)

		evt := <-uiEventCh
		if evt.Err != nil {
			return "", evt.Err
		} else if evt.Type != uii.EventTypeInputReturned {
			return "", fmt.Errorf("wrong EventType: expected %d, got %d", uii.EventTypeInputReturned, evt.Type)
		}
		return evt.Data, nil
	}

	if c.server == "" {
		s, err := handleInput(ui.GetServer)
		if err != nil {
			return err
		}
		c.server = s
	}

	if c.username == "" {
		u, err := handleInput(ui.GetUsername)
		if err != nil {
			return err
		}
		c.username = u
	}

	if c.password == "" {
		p, err := handleInput(ui.GetPassword)
		if err != nil {
			return err
		}
		c.password = p
	}

	// Display information
	//uii.DisplayInformation(server, username, password)

	// Display the list of senders
	listSendersCh := ui.ListSenders(uiEventCh)
	return c.FetchEmails(listSendersCh)
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
