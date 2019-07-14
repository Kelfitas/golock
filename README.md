# GoLock

## Requirements
- go >=1.12

## Install
1. `make build`
1. `make install`

## Usage
- `golock -u user`

## TODO
- ~~hijack keyboard input - ok~~
- ~~handle backspace - ok~~
- ~~handle enter - ok~~
- ~~handle rest of keys - ok~~
- ~~check password on enter (pam_auth): - ok~~
  - ~~on fail    -> clear password - ok~~
  - ~~on success -> quit - ok~~
- ~~add GTK GUI - ok~~
- ~~alert caps lock - ok~~
- ~~get screen size (argv for now) - ok~~
- ~~dynamic user (argv for now) - ok~~
- get rid of deepin dependency
- do everything for each screen
- detect screen changes
- add key bypass list
- fix data races
- tests
