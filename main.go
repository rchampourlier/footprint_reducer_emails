package main

import (
	"log"
	"os"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

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
		log.Fatal(err)
	}

	// List mailboxes
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
		log.Fatal(err)
	}

	for _, m := range mailboxNames {
		mbox, err := c.Select(m, false)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("* " + m)

		if mbox.Messages > 0 {
			// Get all messages
			from := uint32(1)
			to := mbox.Messages
			seqset := new(imap.SeqSet)
			seqset.AddRange(from, to)

			messages := make(chan *imap.Message, 10)
			done = make(chan error, 1)
			go func() {
				done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
			}()

			for msg := range messages {
				log.Println("  - " + msg.Envelope.Subject + " (" + strings.Join(msg.Flags, ",") + ")")
			}

			if err := <-done; err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Println("Done!")
}
