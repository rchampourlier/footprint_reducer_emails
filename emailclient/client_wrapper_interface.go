package emailclient

import (
	"github.com/emersion/go-imap"
)

type ClientWrapper interface {
	Connect(server, email, password string) error
	Logout()
	ListMailboxes() ([]string, error)
	FetchMessages(mailboxName string) ([]*imap.Message, error)
}
