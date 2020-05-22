# User interface

## Program started <a href="program-started"></a>

- ✅ _V0.1_ _UI_ Display a form to enter IMAP server URL and credentials
  - Layout
    - Window in the middle of the screen
    - Description line and input line
    - Validate entry using <ENTER>
    - A single entry at once
- ✅ _V0.1_ _Program_ When confirmed -> [Fetching messages](#fetching-messages)

## List senders <a name="list-senders"></a>

- ⚙️ _V0.1_ _UI_ Display a paginated list of all senders:
  - ✅ Ordered by total size of messages, descending
  - ✅ Support arrows to navigate
    - Up/Down: up/down one line
    - Left/Right: up/down one page
  - ❌ Type <ENTER> to display the current line's sender -> [Show sender](#show-sender)
- _V1.0_ _UI_ Enable going back to the previous screen using <ESC>
- _V1.0_ _UI_ Display instructions on keyboard commands (<ENTER>, arrows, <ESC>) 

## Fetching messages <a name="fetching-messages"></a>

- ✅ _V0.1_ _Program_ Connect to the server
  - _V0.2_ If failed:
    - _UI_ Display an error message
    - _Program_ Go back to [Program started](#program-started)
- ✅ _V0.1_ _Program_ Fetch the messages
- ✅ _V0.1_ _Program_ When messages fetched -> [List senders](#list-senders)
- _V0.2_ _UI_ Display a progress indicator
- _V1.0_ _Program_ When messages fetched -> [Display main menu](#main-menu)

## Show sender <a name="show-sender"></a>

- ❌ _V0.1_ Display a paginated list of all messages for the selected sender:
  - Display the subject
  - Support arrows to navigate
  - Scrolling highlights a message
  - On <ENTER>, display the content of the highlighted message
- ❌ _V0.1_ On <ESC>, go back to the previous screen
- _V0.2_ Display statistics on the sender
  - Total number of emails
  - Total size
  - C02e impact
- _V0.2_ On <BACKSPACE> or <DELETE> -> [Delete sender](#delete-sender)

## Display message <a name="display-message"></a>

- ❌ _V0.1_ Display a full-screen view with the selected message's content
- ❌ _V0.1_ On <ESC>, go back to the previous view
- _V0.2_ On <BACKSPACE> or <DELETE> -> [Delete message](#delete-message)

## Delete message <a name="delete-message"></a>

## Display main menu <a name="display-main-menu"></a>

- _V1.0_ _UI_ Display main menu:
  - Display the email account statistics:
    - Total number of emails
    - Total size
    - C02e generated every year by this data (and an equivalence in food or transport)
  - Display menu items: 
    - _V1.0_ List senders -> [List senders](#list-senders)

## Everywhere

- ✅ _V0.1_ _Feature_ Support exiting on `CTRL-C`

