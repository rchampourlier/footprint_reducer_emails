package main

import (
	"log"
	"os"
	"strconv"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

const MAILBOX_NAME = "[Gmail]/Tous les messages"

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

	messages, err := fetchMessages(c, MAILBOX_NAME)
	if err != nil {
		log.Println("FETCHING MESSAGES ERROR: " + err.Error())
	}
	log.Printf("Done: " + strconv.Itoa(len(messages)) + " messages!\n\n")

	senders := listSenders(messages)
	log.Printf("%d senders\n", len(senders))
	for _, s := range senders {
		log.Printf("  - %s\n", s.Address())
	}
}

func fetchMessages(c *client.Client, mailboxName string) ([]*imap.Message, error) {
	var messages []*imap.Message

	mbox, err := c.Select(mailboxName, false)
	if err != nil {
		log.Println("SELECT MAILBOX ERROR: " + err.Error())
		return nil, err
	}
	if mbox.Messages > 0 {
		messages = make([]*imap.Message, 0)

		// Fetching all messages
		seqset := new(imap.SeqSet)
		seqset.AddRange(1, mbox.Messages)

		fetchedMessages := make(chan *imap.Message, mbox.Messages)
		done := make(chan error, 1)
		go func() {
			done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, fetchedMessages)
		}()

		for msg := range fetchedMessages {
			messages = append(messages, msg)
		}

		if err := <-done; err != nil {
			log.Println("ERROR: " + err.Error())
		}
	}

	return messages, nil
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
