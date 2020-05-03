[![Build Status](https://travis-ci.org/rchampourlier/footprint_reducer_emails.svg?branch=master)](https://travis-ci.org/rchampourlier/footprint_reducer_emails)
[![codecov](https://codecov.io/gh/rchampourlier/footprint_reducer_emails/branch/master/graph/badge.svg)](https://codecov.io/gh/rchampourlier/footprint_reducer_emails)
[![Go Report Card](https://goreportcard.com/badge/github.com/rchampourlier/footprint_reducer_emails)](https://goreportcard.com/report/github.com/rchampourlier/footprint_reducer_emails)

# README

`clean_emails` is a basic Go program to help you clean you email history. It should be compatible with any IMAP server.

## STATUS

For now it's a simple script demonstrating we can connect to an IMAP server and fetch some emails.

## HOW TO USE

```
EMAIL=REPLACEME PASSWORD=REPLACEME go run main.go
```

## DESIGN DECISIONS

### Command-line with parameters or TUI?

I chose to make a TUI (Terminal User Interface). 

Since the goal of the app is to reduce the environmental footprint of your email inboxes, it made sense to limit the footprint of the app itself. In particular, the bandwidth required by the app to do its job.

Having an app that could be run with parameters to do its job would have needed either to have a persistent storage to use it between launches of the app, making it more complex.

Instead, thinking of the use of the app as a session where the user connects and fetch data from the email server once, then take decisions and gives order (e.g. delete emails from this sender), made sense. It limited the data to fetch and enabled to simply store the data in memory during a "session", without persisting data once the program is ended.
