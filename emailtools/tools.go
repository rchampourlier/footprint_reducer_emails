package emailtools

import (
	"sort"
	"time"

	"github.com/emersion/go-imap"
)

// SenderStat is a structure containing basic statistics on a
// sender.
type SenderStat struct {
	Sender            *imap.Address
	MessagesCount     uint
	LatestMessageDate time.Time
	TotalSize         uint32
}

// ListSenders returns the list of all senders present in the passed
// `messages`.
//
// The result is a slice of unique `*imap.Address`. All senders are
// included, even when several senders are present in a single messsage.
func ListSenders(messages []*imap.Message) []*imap.Address {
	uniqueSenders := make(map[string]bool)
	senders := make([]*imap.Address, 0)

	for _, m := range messages {
		for _, msgSender := range m.Envelope.Sender {
			if uniqueSenders[msgSender.Address()] != true {
				uniqueSenders[msgSender.Address()] = true
				senders = append(senders, msgSender)
			}
		}
	}

	return senders
}

// StatsOnSenders returns a slice of *SenderStat, with the statistics for each
// sender in the given list of messages.
func StatsOnSenders(messages []*imap.Message) []*SenderStat {
	statsMap := make(map[string]*SenderStat)
	stats := make([]*SenderStat, 0)

	for _, m := range messages {
		for _, msgSender := range m.Envelope.Sender {
			if statsMap[msgSender.Address()] != nil {
				stat := statsMap[msgSender.Address()]
				stat.MessagesCount++
				if stat.LatestMessageDate.Before(m.Envelope.Date) {
					stat.LatestMessageDate = m.Envelope.Date
				}
				stat.TotalSize += m.Size
			} else {
				newStat := SenderStat{
					Sender:            msgSender,
					MessagesCount:     1,
					LatestMessageDate: m.Envelope.Date,
					TotalSize:         m.Size,
				}
				statsMap[msgSender.Address()] = &newStat
				stats = append(stats, &newStat)
			}
		}
	}

	return stats
}

// SortSendersStatBySize sorts the passed slice of `*SenderStat` on
// its `TotalSize` field, descending.
//
// The sort is performed in place.
func SortSendersStatBySize(s []*SenderStat) {
	sort.Slice(
		s,
		func(i, j int) bool {
			return s[i].TotalSize > s[j].TotalSize
		},
	)
}
