//+build !test

package main

import (
	"log"
	"os"
	"strconv"

	"footprint_reducer_emails/emailclient"
	"footprint_reducer_emails/emailtools"
)

const mailboxName = "[Gmail]/Tous les messages"

func main() {
	c, err := emailclient.ConnectAndLogin(
		os.Getenv("SERVER"),
		os.Getenv("EMAIL"),
		os.Getenv("PASSWORD"),
	)
	// Don't forget to logout
	defer c.Logout()

	messages, err := c.FetchMessages(mailboxName)
	if err != nil {
		log.Println("FETCHING MESSAGES ERROR: " + err.Error())
	}
	log.Printf("Done: " + strconv.Itoa(len(messages)) + " messages!\n\n")

	senders := emailtools.ListSenders(messages)
	log.Printf("%d senders\n\n", len(senders))

	stats := emailtools.StatsOnSenders(messages)
	emailtools.SortSendersStatBySize(stats)

	var totalMailboxSize uint32
	for _, stat := range stats {
		totalMailboxSize += stat.TotalSize
		log.Printf("  - %s: %d messages for %d MB, latest message on %s\n", stat.Sender.Address(), stat.MessagesCount, stat.TotalSize/1024^2, stat.LatestMessageDate)
	}

	log.Printf("\nTotal mailbox size: %d MB\n", totalMailboxSize/1024^2)
}
