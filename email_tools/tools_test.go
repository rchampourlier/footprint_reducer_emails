package email_tools_test

import (
	"footprint_reducer_emails/email_tools"

	"testing"
	"time"

	"github.com/emersion/go-imap"
)

// Given a slice of emails, ListSenders should return a slice of
// unique IMAP addresses of the messages.
//
// - If the same sender is present in several emails, it should
//   be present only once in the returned slice.
// - If a message has several senders, all senders should be
//   present in the returned slice.
func TestListSenders(t *testing.T) {
	// Test uniqueness in different and same message

	// Create 2 messages with the same sender
	// Assert there is only 1 sender
	messages := []*imap.Message{
		createMessageWithSender([]string{"sender1"}),
		createMessageWithSender([]string{"sender1"}),
	}
	senders := email_tools.ListSenders(messages)
	if len(senders) != 1 {
		t.Fatalf("Expected sender uniqueness from different messages (expected 1, found %d)\n", len(senders))
	}

	// Create 1 message with the same sender twice
	// Assert there is only 1 sender
	messages = []*imap.Message{
		createMessageWithSender([]string{"sender1", "sender1"}),
	}
	senders = email_tools.ListSenders(messages)
	if len(senders) != 1 {
		t.Fatalf("Expected sender uniqueness in single message (expected 1, found %d)\n", len(senders))
	}

	// Test completeness
	// Create 2 messages with each 2 different senders (all differents)
	// Assert there are 4 senders
	messages = []*imap.Message{
		createMessageWithSender([]string{"sender1", "sender2"}),
		createMessageWithSender([]string{"sender3", "sender4"}),
	}
	senders = email_tools.ListSenders(messages)
	if len(senders) != 4 {
		t.Fatalf("Expected sender completeness (expected 4, got %d)\n", len(senders))
	}
}

// Assertions:
// - Test all fields of `SenderStat`:
//   - total size,
//   - latest message date,
//   - messages count.
// - Check uniqueness and completeness.
func TestStatsOnSenders(t *testing.T) {
	date1 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC)
	messages := []*imap.Message{
		createMessage(
			[]string{"sender1", "sender2"},
			1,
			date1,
		),
		createMessage(
			[]string{"sender1"},
			10,
			date2,
		),
		createMessage(
			[]string{"sender3"},
			100,
			date3,
		),
	}
	SenderStats := email_tools.StatsOnSenders(messages)

	// Check uniqueness and completeness
	if len(SenderStats) != 3 {
		t.Fatalf("Expects uniqueness and completeness of senders (expects 3, got %d)\n", len(SenderStats))
	}

	// Check fields
	for _, stat := range SenderStats {
		switch stat.Sender.MailboxName {

		case "sender1":
			// Expects msg count = 2, size = 11, date2
			if stat.MessagesCount != 2 {
				t.Fatalf("Expects message count of 2 for `sender1`, got %d\n", stat.MessagesCount)
			}
			if stat.TotalSize != 11 {
				t.Fatalf("Expects total size of 11 for `sender1`, got %d\n", stat.TotalSize)
			}
			if stat.LatestMessageDate != date2 {
				t.Fatalf("Expects latest message date for `sender1` to be %x, got %x\n", date2, stat.LatestMessageDate)
			}

		case "sender2":
			// Expects msg count = 1, size = 1, date1
			if stat.MessagesCount != 1 {
				t.Fatalf("Expects message count of 1 for `sender2`, got %d\n", stat.MessagesCount)
			}
			if stat.TotalSize != 1 {
				t.Fatalf("Expects total size of 1 for `sender2`, got %d\n", stat.TotalSize)
			}
			if stat.LatestMessageDate != date1 {
				t.Fatalf("Expects latest message date for `sender2` to be %x, got %x\n", date2, stat.LatestMessageDate)
			}
		case "sender3":
			// Expects msg count 1, size 100, date3
			if stat.MessagesCount != 1 {
				t.Fatalf("Expects message count of 1 for `sender3`, got %d\n", stat.MessagesCount)
			}
			if stat.TotalSize != 100 {
				t.Fatalf("Expects total size of 100 for `sender3`, got %d\n", stat.TotalSize)
			}
			if stat.LatestMessageDate != date3 {
				t.Fatalf("Expects latest message date for `sender3` to be %x, got %x\n", date3, stat.LatestMessageDate)
			}
		}
	}
}

func TestSortSendersStatBySize(t *testing.T) {
	SenderStats := []*email_tools.SenderStat{
		{
			Sender:    createAddress("sender1"),
			TotalSize: 100,
		},
		{
			Sender:    createAddress("sender2"),
			TotalSize: 1,
		},
	}
	email_tools.SortSendersStatBySize(SenderStats)
	if SenderStats[0].Sender.MailboxName != "sender1" {
		t.Fatalf("Expected the slice to be sorted by total size, descending")
	}
}

func createMessageWithSender(senderNames []string) *imap.Message {
	senders := createSenders(senderNames)
	return &imap.Message{
		Envelope: &imap.Envelope{
			Sender: senders,
		},
	}
}

func createMessage(senderNames []string, size uint32, date time.Time) *imap.Message {
	senders := createSenders(senderNames)
	return &imap.Message{
		Envelope: &imap.Envelope{
			Sender: senders,
			Date:   date,
		},
		Size: size,
	}
}

func createSenders(senderNames []string) []*imap.Address {
	senders := make([]*imap.Address, len(senderNames))
	for i, name := range senderNames {
		senders[i] = createAddress(name)
	}
	return senders
}

func createAddress(senderName string) *imap.Address {
	return &imap.Address{
		MailboxName: senderName,
		HostName:    "host",
	}
}
