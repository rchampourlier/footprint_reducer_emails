package email_client

import (
	"github.com/emersion/go-imap"
)

type ImapClient interface {
	Select(name string, readOnly bool) (*imap.MailboxStatus, error)
	Fetch(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error
}
