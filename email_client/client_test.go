package email_client_test

import (
	"footprint_reducer_emails/email_client"
	"testing"

	"github.com/emersion/go-imap"
)

// Expected to return a slice of the messages for the specified
// mailbox.
func TestFetchMessages(t *testing.T) {
	clientMock := email_client.NewMockImapClient(t)

	mailboxStatus := &(imap.MailboxStatus{
		Messages: 2,
	})
	clientMock.ExpectSelect().
		WillRespondWithMailboxStatus(mailboxStatus)

	messages := []*imap.Message{
		&imap.Message{
			Envelope: &imap.Envelope{
				Sender: []*imap.Address{
					{
						MailboxName: "sender1",
						HostName:    "host1",
					},
				},
			},
			Size: 100,
		},
		&imap.Message{
			Envelope: &imap.Envelope{
				Sender: []*imap.Address{
					{
						MailboxName: "sender2",
						HostName:    "host2",
					},
				},
			},
			Size: 200,
		},
	}
	clientMock.ExpectFetch().
		WillRespondWith(nil).
		WillSend(messages)

	fetchedMessages, err := email_client.FetchMessages(clientMock, "mailbox#0")
	if err != nil {
		t.Fatalf("FetchMessages returned an error: %s\n", err)
	}
	//log.Printf("in test: %d\n", len(fetchedMessages))
	if len(fetchedMessages) != 2 {
		t.Fatalf("Expected FetchMessages to return 2 messages, got %d\n", len(fetchedMessages))
	}

	// expects to select mailbox
	// expects to fetch messages for a sequence set from 1 to the
	//   total number of messages
	// expects to fetch envelope and size for all messages
	// expects to return a slice of `*imap.Message`
}

func TestFetchMessagesSelectError(t *testing.T) {
	// expects to select mailbox
	// triggers error
	// expects to return nil, error
}

func TestFetchMessagesFetchError(t *testing.T) {
	// expects to select mailbox
	// expects to ignore the error, log a message
	// expects to return a slice of `*imap.Message` without the one
	//   in error
}
