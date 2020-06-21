package emailclient

import (
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// ImapClientWrapper wraps an Imap email client.
type ImapClientWrapper struct {
	c ImapClient
}

// NewImapClientWrapper returns a `ImapClientWrapper` struct.
func NewImapClientWrapper(c ImapClient) *ImapClientWrapper {
	return &ImapClientWrapper{c}
}

// Connect connects the client to the server, then logins
// if the connection succeeded.
//
// It stores the reference to the `imap.Client` it wraps.
// Returns an error if the connection fails.
func (c *ImapClientWrapper) Connect(server, email, password string) error {
	imapClient, err := client.DialTLS(server, nil)
	if err != nil {
		return err
	}
	c.c = imapClient
	if err = imapClient.Login(email, password); err != nil {
		log.Fatalln("LOGIN ERROR: " + err.Error())
	}
	return err
}

// Logout logs the client ouf of the server.
// Should be called in a `defer` after `Connect`.
func (c *ImapClientWrapper) Logout() {
	c.c.Logout()
}

// ListMailboxes fetches the list of mailboxes available on the server
// and return a slice of their names or an error.
func (c *ImapClientWrapper) ListMailboxes() ([]string, error) {
	var err error
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.c.List("", "*", mailboxes)
	}()

	mailboxNames := make([]string, 0)
	for m := range mailboxes {
		mailboxNames = append(mailboxNames, m.Name)
	}

	if err = <-done; err != nil {
		log.Println("LIST MAILBOX ERROR: " + err.Error())
	}
	return mailboxNames, err
}

// FetchMessages fetches all messages using the specified `go-imap/client.ClientWrapper`,
// from the specified mailbox.
// Returns a slice of `*imap.Message` or an error.
func (c *ImapClientWrapper) FetchMessages(mailboxName string) ([]*imap.Message, error) {
	var messages []*imap.Message

	mbox, err := c.c.Select(mailboxName, false)
	if err != nil {
		log.Println("SELECT MAILBOX ERROR: " + err.Error())
		return nil, err
	}
	if mbox.Messages > 0 {
		messages = make([]*imap.Message, 0)

		// Fetching all messages
		seqset := new(imap.SeqSet)
		seqset.AddRange(1, mbox.Messages)

		fetchedMessages := make(chan *imap.Message, mbox.Messages)
		done := make(chan error, 1)
		go func() {
			done <- c.c.Fetch(
				seqset,
				[]imap.FetchItem{
					imap.FetchEnvelope,
					imap.FetchRFC822Size,
				},
				fetchedMessages)
		}()

		for msg := range fetchedMessages {
			messages = append(messages, msg)
		}

		if err := <-done; err != nil {
			log.Println("ERROR: " + err.Error())
		}
	}

	return messages, nil
}
