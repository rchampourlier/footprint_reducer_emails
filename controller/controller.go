package controller

import (
	"fmt"

	"footprint_reducer_emails/emailclient"
	"footprint_reducer_emails/emailtools"
	uii "footprint_reducer_emails/ui"

	"github.com/emersion/go-imap"
)

// Controller represents a controller and stored the reference
// to the UI and the state of the program execution.
type Controller struct {
	ui uii.UI

	// Data
	server   string
	username string
	password string
	messages []*imap.Message

	// Calculated data
	senderStats []*emailtools.SenderStat
}

// NewController returns a new controller with the specified UI
func NewController(i uii.UI) *Controller {
	return NewControllerWithCredentials(i, "", "", "")
}

// NewControllerWithCredentials initializes a new controller with
// the specified UI and server  and credentials.
func NewControllerWithCredentials(i uii.UI, server, username, password string) *Controller {
	msgs := make([]*imap.Message, 0)
	ss := make([]*emailtools.SenderStat, 0)
	return &Controller{i, server, username, password, msgs, ss}
}

// Run executes the program.
func (c *Controller) Run() error {
	ui := c.ui
	uiEventCh := make(chan uii.Event, 0)
	defer close(uiEventCh)

	handleInputReturned := func() (string, error) {
		evt := <-uiEventCh
		if evt.Err != nil {
			return "", evt.Err
		} else if evt.Type != uii.EventTypeStringInputReturned {
			return "", fmt.Errorf("wrong EventType: expected %d, got %d", uii.EventTypeStringInputReturned, evt.Type)
		}
		return evt.Data.(string), nil
	}

	if c.server == "" {
		if err := ui.StringInput("Enter the server URL and port (e.g. imap.gmail.com:993):", uiEventCh); err != nil {
			return err
		}
		data, err := handleInputReturned()
		if err != nil {
			return err
		}
		c.server = data
	}

	if c.username == "" {
		if err := ui.StringInput("Enter your IMAP username (generally your email address):", uiEventCh); err != nil {
			return err
		}
		data, err := handleInputReturned()
		if err != nil {
			return err
		}
		c.username = data
	}

	if c.password == "" {
		if err := ui.StringWithMaskInput("Enter your IMAP password:", '•', uiEventCh); err != nil {
			return err
		}
		data, err := handleInputReturned()
		if err != nil {
			return err
		}
		c.password = data
	}

	// Fetch messages
	if err := c.fetchMessages(); err != nil {
		return err
	}

	// Calculate senderStats
	c.senderStats = emailtools.StatsOnSenders(c.messages)
	emailtools.SortSendersStatBySize(c.senderStats)

	// Display senders
	senderLines := make([]string, 0)
	for i, stat := range c.senderStats {
		line := fmt.Sprintf("%04d | %s | %d messages | %d MB | %s", i, stat.Sender.Address(), stat.MessagesCount, stat.TotalSize/1024^2, stat.LatestMessageDate)
		senderLines = append(senderLines, line)
	}
	ui.List(senderLines, uiEventCh)

	// Waiting for an event on the list of senders
	evt := <-uiEventCh
	if evt.Err != nil {
		return evt.Err
	} else if evt.Type != uii.EventTypeItemSelected {
		return fmt.Errorf("invalid ui.EventType (expected %d, got %d)", uii.EventTypeItemSelected, evt.Type)
	}

	selectedSenderIndex := evt.Data.(int)
	selectedSender := c.senderStats[selectedSenderIndex].Sender

	// Display messages for the selected sender
	messageLines := make([]string, 0)
	for _, msg := range c.messagesForSenderAddress(selectedSender) {
		messageLines = append(messageLines, msg)
	}
	ui.List(messageLines, uiEventCh)

	return nil
}

func (c *Controller) fetchMessages() error {
	client, err := emailclient.ConnectAndLogin(c.server, c.username, c.password)
	if err != nil {
		return fmt.Errorf("failed to connect to IMAP server: %w", err)
	}

	// Don't forget to logout
	defer client.Logout()

	// TODO remove this constant
	const mailboxName = "[Gmail]/Tous les messages"

	messages, err := client.FetchMessages(mailboxName)
	if err != nil {
		return err
	}
	c.messages = messages
	return nil
}

// messagesForSenderAddress returns a slice of strings where each line represent
// a message of the specified sender.
// Messages must have been fetched before with `fetchMessages`.
func (c *Controller) messagesForSenderAddress(sender *imap.Address) []string {
	msgs := emailtools.MessagesForSenderAddress(sender, c.messages)
	emailtools.SortMessagesBySize(msgs)

	lines := make([]string, 0)
	for i, msg := range msgs {
		line := fmt.Sprintf("%04d | %.0f MB | %s", i, float32(msg.Size/1024^2), msg.Envelope.Subject)
		lines = append(lines, line)
	}

	return lines
}
