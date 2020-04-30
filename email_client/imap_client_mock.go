package email_client

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/emersion/go-imap"
)

// MockImapClient is the base struct to build a Mock.
//
// This pattern for mocking is inspired from
// [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock).
type MockImapClient struct {
	*testing.T
	expectations []Expectation
	mutex        sync.Mutex
}

// Expectation is a specific interface for structs representing
// expectations for the mock. They implement a `Describe` method
// that can be used by the mock to display when there is a
// mismatch between the expected call and the call it received.
type Expectation interface {
	Describe() string
}

// NewMockImapClient returns a new `MockImapClient` with a default
// behaviour.
func NewMockImapClient(t *testing.T) *MockImapClient {
	return &MockImapClient{
		T:     t,
		mutex: sync.Mutex{},
	}
}

// List
func (m *MockImapClient) List(ref, name string, ch chan *imap.MailboxInfo) error {
	msg := "ImapClientMock.List not implemented"
	log.Fatalln(msg)
	return fmt.Errorf(msg)
}

// Fetch
func (m *MockImapClient) Fetch(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
	e := m.popExpectation()
	if e == nil {
		m.Errorf("mock received `Fetch` but no expectation was set")
	}
	ee, ok := e.(*ExpectedFetch)
	if !ok {
		m.Errorf("mock received `Fetch` but was expecting `%s`\n", e.Describe())
	}
	// Send ee.messages over ch
	for _, m := range ee.messages {
		ch <- m
	}
	close(ch)
	return ee.err
}

// Select
func (m *MockImapClient) Select(name string, readOnly bool) (*imap.MailboxStatus, error) {
	e := m.popExpectation()
	if e == nil {
		m.Errorf("mock received `Select` but no expectation was set")
	}
	ee, ok := e.(*ExpectedSelect)
	if !ok {
		m.Errorf("mock received `Select` but was expecting `%s`\n", e.Describe())
	}
	return ee.mailboxStatus, ee.err
}

// ============
// Expectations
// ============

// Fetch
// ----------

// ExpectedFetch is an expectation for `Fetch`
//
// Use `With...` and `Will...` methods on the returned
// `ExpectedReplaceIssueStateAndEvents` expectation to
// specify expected arguments and return value.
type ExpectedFetch struct {
	err      error
	messages []*imap.Message
	// add the expectation's parameters to be checked when the expected
	// method is called
}

// ExpectFetch indicates the mock should expect a call to
// `Fetch` with the specified arguments.
func (m *MockImapClient) ExpectFetch() *ExpectedFetch {
	e := ExpectedFetch{}
	m.expectations = append(m.expectations, &e)
	return &e
}

// Describe describes the `Fetch` expectation
func (e *ExpectedFetch) Describe() string {
	return fmt.Sprintf("Fetch with args...")
}

// WillRespondWith indicates `ExpectedFetch`
// expectation should return the specified value when
// called.
// Returns itself so calls to `Will...` may be chained.
func (e *ExpectedFetch) WillRespondWith(err error) *ExpectedFetch {
	e.err = err
	return e
}

func (e *ExpectedFetch) WillSend(messages []*imap.Message) *ExpectedFetch {
	e.messages = messages
	return e
}

// ExpectedSelect is an expectation for `Select`
//
// Use `With...` and `Will...` methods on the returned
// `ExpectedReplaceIssueStateAndEvents` expectation to
// specify expected arguments and return value.
type ExpectedSelect struct {
	mailboxStatus *imap.MailboxStatus
	err           error
	// add the expectation's parameters to be checked when the expected
	// method is called
}

// ExpectSelect indicates the mock should expect a call to
// `Select` with the specified arguments.
func (m *MockImapClient) ExpectSelect() *ExpectedSelect {
	e := ExpectedSelect{}
	m.expectations = append(m.expectations, &e)
	return &e
}

// Describe describes the `Select` expectation
func (e *ExpectedSelect) Describe() string {
	return fmt.Sprintf("Select with args...")
}

// WillRespondWithMailboxStatus will make a call to
// `Select` on the mock return the specified value and
// a nil error.
func (e *ExpectedSelect) WillRespondWithMailboxStatus(mailboxStatus *imap.MailboxStatus) *ExpectedSelect {
	e.mailboxStatus = mailboxStatus
	return e
}

// Other
// -----

func (m *MockImapClient) popExpectation() Expectation {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if len(m.expectations) == 0 {
		return nil
	}
	e := m.expectations[0]
	m.expectations = m.expectations[1:]
	return e
}
