package emailclient

import (
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// Client represents a client supporting the connection to
// an Imap server.
type Client struct {
	c ImapClient
}

// NewClient returns a pointer to a `Client` struct
// with the specified `ImapClient`.
func NewClient(c ImapClient) *Client {
	return &Client{c}
}

// ConnectAndLogin connects the client to the server, then logins
// if the connection succeeded.
//
// It stores the reference to the `imap.Client` is uses.
// Returns an error if the connection fails.
func ConnectAndLogin(server, email, password string) (*Client, error) {
	imapClient, err := client.DialTLS(server, nil)
	c := Client{imapClient}
	if err != nil {
		log.Fatal(err)
		return &c, err
	}
	if err = c.c.Login(email, password); err != nil {
		log.Fatalln("LOGIN ERROR: " + err.Error())
	}
	return &c, err
}

// Logout logs the client ouf of the server.
// Should be called in a `defer` after `Connect`.
func (c *Client) Logout() {
	c.c.Logout()
}

// ListMailboxes fetches the list of mailboxes available on the server
// and return a slice of their names or an error.
func (c *Client) ListMailboxes() ([]string, error) {
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

// FetchMessages fetches all messages using the specified `go-imap/client.Client`,
// from the specified mailbox.
// Returns a slice of `*imap.Message` or an error.
func (c *Client) FetchMessages(mailboxName string) ([]*imap.Message, error) {
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
