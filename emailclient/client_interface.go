package emailclient

import (
	"github.com/emersion/go-imap"
)

// ImapClient defines an interface for `imap.Client` so it may get
// another implementation, in particular a mock (`ImapClientMock`)
// we may use for testing.
type ImapClient interface {
	List(ref, name string, ch chan *imap.MailboxInfo) error
	Select(name string, readOnly bool) (*imap.MailboxStatus, error)
	Fetch(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error
}
