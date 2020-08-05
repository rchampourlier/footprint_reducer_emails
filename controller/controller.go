package controller

import (
	"fmt"
	"log"
	"os"

	"footprint_reducer_emails/emailclient"
	"footprint_reducer_emails/emailtools"
	uii "footprint_reducer_emails/ui"

	"github.com/emersion/go-imap"
)

// TODO remove this constant
const mailboxName = "[Gmail]/Tous les messages"

// Controller represents a controller and stored the reference
// to the UI and the state of the program execution.
type Controller struct {
	w  emailclient.ClientWrapper
	ui uii.UI

	// Data
	server   string
	email    string
	password string
	messages []*imap.Message

	// Calculated data
	senderStats []*emailtools.SenderStat
}

// NewController returns a new controller with the specified
// Imap client wrapper and UI.
//
// It attempts to retrieve the server URL and credentials from
// `SERVER`, `EMAIL` and `PASSWORD` environment variables.
func NewController(w emailclient.ClientWrapper, i uii.UI) *Controller {
	server := os.Getenv("SERVER")
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	return NewControllerWithCredentials(w, i, server, email, password)

}

// NewControllerWithCredentials initializes a new controller with
// the specified UI and server  and credentials.
func NewControllerWithCredentials(w emailclient.ClientWrapper, i uii.UI, server, email, password string) *Controller {
	msgs := make([]*imap.Message, 0)
	ss := make([]*emailtools.SenderStat, 0)
	return &Controller{w, i, server, email, password, msgs, ss}
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

	if c.email == "" {
		if err := ui.StringInput("Enter your IMAP email (generally your email address):", uiEventCh); err != nil {
			return err
		}
		data, err := handleInputReturned()
		if err != nil {
			return err
		}
		c.email = data
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

	// Connect to Imap server
	if err := c.w.Connect(c.server, c.email, c.password); err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	// Don't forget to logout
	defer c.w.Logout()

	// Fetch messages
	if err := c.fetchMessages(); err != nil {
		return err
	}

	if len(c.messages) == 0 {
		ui.Information(
			"Messages fetched",
			"Done, 0 messages\nExit with CTRL-C.",
		)
		return nil
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
	sortedMessagesForSelectedSender := c.messagesForSenderAddressSortedBySize(selectedSender)
	for i, msg := range sortedMessagesForSelectedSender {
		messageLines = append(messageLines, listItemFromMessage(i, msg))
	}
	ui.List(messageLines, uiEventCh)

	// Waiting for an event on the list of messages
	evt = <-uiEventCh
	if evt.Err != nil {
		return evt.Err
	} else if evt.Type != uii.EventTypeItemSelected {
		return fmt.Errorf("invalid ui.EventType (expected %d, got %d)", uii.EventTypeItemSelected, evt.Type)
	}

	selectedMessageIndex := evt.Data.(int)
	selectedMessage := sortedMessagesForSelectedSender[selectedMessageIndex]

	// Fetch selected message body
	selectedMessage, err := c.fetchMessageText(selectedMessage)
	if err != nil {
		logger().Printf("error fetching message body: %x\n", err)
	}

	bodyString := ""
	// Display selected message
	for _, sectionLiteral := range selectedMessage.Body {
		literalBytes := make([]byte, sectionLiteral.Len())
		_, err := sectionLiteral.Read(literalBytes)

		// Stripping \r from the literal since it breaks the display
		// in gocui's views.
		// TODO: move this to the ui library instead.
		crStrippedLiteral := make([]byte, 0)
		for i := 0; i < len(literalBytes); i++ {
			if literalBytes[i] != '\r' {
				crStrippedLiteral = append(crStrippedLiteral, literalBytes[i])
			}
		}

		if err != nil {
			logger().Println("error reading message body")
		} else {
			bodyString += fmt.Sprintf("%s\n---\n", crStrippedLiteral)
		}
	}
	ui.Page(selectedMessage.Envelope.Subject, bodyString)

	return nil
}

func (c *Controller) fetchMessages() error {
	messages, err := c.w.FetchMessages(mailboxName)
	if err != nil {
		return err
	}
	c.messages = messages

	return nil
}

// fetchMessageText fetches the text parts of the messages
// displayable on a text-only client.
// The fetched body is added to the `imap.Message` struct.
func (c *Controller) fetchMessageText(m *imap.Message) (*imap.Message, error) {
	return c.w.FetchMessageText(mailboxName, m.Uid)
}

// messagesForSenderAddress returns a slice of strings where each line represent
// a message of the specified sender.
// Messages must have been fetched before with `fetchMessages`.
func (c *Controller) messagesForSenderAddressSortedBySize(sender *imap.Address) []*imap.Message {
	msgs := emailtools.MessagesForSenderAddress(sender, c.messages)
	emailtools.SortMessagesBySize(msgs)
	return msgs
}

func listItemFromMessage(i int, m *imap.Message) string {
	return fmt.Sprintf("%04d | %.0f MB | %s", i, float32(m.Size/1024^2), m.Envelope.Subject)
}

func logger() *log.Logger {
	f, err := os.Create("./log.txt")
	if err != nil {
		panic(err)
	}
	logger := log.New(f, "", 0)
	return logger
}
