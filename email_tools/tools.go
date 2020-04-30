package email_tools

import (
	"sort"
	"time"

	"github.com/emersion/go-imap"
)

type senderStat struct {
	Sender            *imap.Address
	MessagesCount     uint
	LatestMessageDate time.Time
	TotalSize         uint32
}

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

func SortSendersStatBySize(s []*senderStat) {
	sort.Slice(
		s,
		func(i, j int) bool {
			return s[i].TotalSize > s[j].TotalSize
		},
	)
}
func StatsOnSenders(messages []*imap.Message) []*senderStat {
	statsMap := make(map[string]*senderStat)
	stats := make([]*senderStat, 0)

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
				newStat := senderStat{
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
