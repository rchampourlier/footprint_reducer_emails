package emailclient_test

import (
	"footprint_reducer_emails/emailclient"
	"testing"

	"github.com/emersion/go-imap"
)

// Expected to return a slice of the messages for the specified
// mailbox.
func TestFetchMessages(t *testing.T) {
	mockImapClient := emailclient.NewMockImapClient(t)
	testedClient := emailclient.NewImapClientWrapper(mockImapClient)

	mailboxStatus := &(imap.MailboxStatus{
		Messages: 2,
	})
	mockImapClient.ExpectSelect().
		WillRespondWithMailboxStatus(mailboxStatus)

	mockImapClient.ExpectFetch().
		WillRespondWith(nil).
		WillSend(fixtureMessages())

		// TODO: here we are testing the client wrapper!
	fetchedMessages, err := testedClient.FetchMessages("mailbox#0")
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

func fixtureMessages() []*imap.Message {
	return []*imap.Message{
		{
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
		{
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
}
