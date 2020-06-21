package emailclient

import (
	"log"
	"os"

	"github.com/emersion/go-imap"
)

type MockClientWrapper struct {
	logger *log.Logger
}

// NewMockClientWrapper returns a new `MockClientWrapper`.
// It gets an initialized logger that will write logs to ./log.txt.
func NewMockClientWrapper() *MockClientWrapper {
	f, err := os.Create("./log.txt")
	if err != nil {
		panic(err)
	}
	logger := log.New(f, "", 0)
	return &MockClientWrapper{logger}
}

// Connect only writes a log.
func (w *MockClientWrapper) Connect(server, email, password string) error {
	w.logger.Println("Connect")
	return nil
}

// Logout only writes a log.
func (w *MockClientWrapper) Logout() {
	w.logger.Println("Logout")
}

func (w *MockClientWrapper) ListMailboxes() ([]string, error) {
	mailboxes := make([]string, 0)
	w.logger.Println("ListMailboxes")
	return mailboxes, nil
}

func (w *MockClientWrapper) FetchMessages(mailboxName string) ([]*imap.Message, error) {
	messages := FixtureMessages()
	w.logger.Println("FetchMessages")
	return messages, nil
}

// FixtureMessages generates a slice of `*imap.Message`
// to be used as fixtures for test or to test an application
// in development.
func FixtureMessages() []*imap.Message {
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
