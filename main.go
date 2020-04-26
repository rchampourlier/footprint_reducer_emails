package main

import (
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"

	"footprint_reducer_emails/email_client"
)

const mailboxName = "[Gmail]/Tous les messages"

func main() {
	// Connect to server
	server := os.Getenv("SERVER")
	c, err := client.DialTLS(server, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Don't forget to logout
	defer c.Logout()

	// Login
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	if err := c.Login(email, password); err != nil {
		log.Fatalln("LOGIN ERROR: " + err.Error())
	}

	messages, err := email_client.FetchMessages(c, mailboxName)
	if err != nil {
		log.Println("FETCHING MESSAGES ERROR: " + err.Error())
	}
	log.Printf("Done: " + strconv.Itoa(len(messages)) + " messages!\n\n")

	senders := listSenders(messages)
	log.Printf("%d senders\n\n", len(senders))

	stats := statsOnSenders(messages)
	sortSendersStatBySize(stats)

	var totalMailboxSize uint32
	for _, stat := range stats {
		totalMailboxSize += stat.TotalSize
		log.Printf("  - %s: %d messages for %d MB, latest message on %s\n", stat.Sender.Address(), stat.MessagesCount, stat.TotalSize/1024^2, stat.LatestMessageDate)
	}

	log.Printf("\nTotal mailbox size: %d MB\n", totalMailboxSize/1024^2)
}

func listMailboxes(c *client.Client) []string {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	mailboxNames := make([]string, 0)
	for m := range mailboxes {
		mailboxNames = append(mailboxNames, m.Name)
	}

	if err := <-done; err != nil {
		log.Fatalln("LIST MAILBOX ERROR: " + err.Error())
	}
	return mailboxNames
}

type senderStat struct {
	Sender            *imap.Address
	MessagesCount     uint
	LatestMessageDate time.Time
	TotalSize         uint32
}

func statsOnSenders(messages []*imap.Message) []*senderStat {
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

func sortSendersStatBySize(s []*senderStat) {
	sort.Slice(
		s,
		func(i, j int) bool {
			return s[i].TotalSize > s[j].TotalSize
		},
	)
}

func listSenders(messages []*imap.Message) []*imap.Address {
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
