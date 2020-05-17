# User interface

## Program started <a href="program-started"></a> ✅

- _UI_ Display a form to enter IMAP server URL and credentials
  - Layout
    - Window in the middle of the screen
    - Description line and input line
    - Validate entry using <ENTER>
    - A single entry at once
- _Program_ When confirmed -> [Fetching messages](#fetching-messages)

## Fetching messages <a name="fetching-messages"></a>

- _Program_ Connect to the server ✅
  - If failed:
    - _UI_ Display an error message ❌
    - _Program_ Go back to [Program started](#program-started) ❌
- _Program_ Fetch the messages ✅
- _UI_ Display a progress indicator _V0.X_ _LATER_
- _Program_ When messages fetched -> [Display main menu](#main-menu) ✅
  
## Display main menu <a name="display-main-menu"></a>

- _UI_ Display main menu:
  - Display the email account statistics:
    - Total number of emails _V0.1_
    - Total size _V0.1_
    - C02e generated every year by this data (and an equivalence in food or transport) _V0.X_ _LATER_
  - Display menu items: 
    - List senders _V0.1_ -> [List senders](#list-senders)

## List senders <a name="list-senders"></a>

- _UI_ Display a paginated list of all senders:
  - Ordered by total size of messages, descending ✅
  - Support arrows to navigate: ✅
    - Left/Right: up/down one page
    - Up/Down: up/down one line
  - Type <ENTER> to display the current line's sender -> [Show sender](#show-sender) ❌
- _UI_ Enable going back to the previous screen using <ESC> ❌
- _UI_ Display instructions on keyboard commands (<ENTER>, arrows, <ESC>) ❌

## Show sender <a name="show-sender"></a> ❌

## Everywhere

- _Feature_ Support exiting on `CTRL-C` ✅

