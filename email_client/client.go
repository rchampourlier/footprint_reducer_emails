package email_client

import (
	"log"

	"github.com/emersion/go-imap"
)

// Fetches the list of mailboxes available on the server
// and return a slice of their names or an error.
func ListMailboxes(c ImapClient) ([]string, error) {
	var err error
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
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

// Fetches all messages using the specified `go-imap/client.Client`,
// from the specified mailbox.
// Returns a slice of `*imap.Message` or an error.
func FetchMessages(c ImapClient, mailboxName string) ([]*imap.Message, error) {
	var messages []*imap.Message

	mbox, err := c.Select(mailboxName, false)
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
			done <- c.Fetch(
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
