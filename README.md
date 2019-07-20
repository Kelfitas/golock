# GoLock

## Requirements
- go >=1.12

## Install
1. `make build`
1. `make install`

## Usage
- `golock`

### Flags
<table>
    <thead>
        <tr>
            <th>flag</th>
            <th>type</th>
            <th>desc</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>-a</td>
            <td>float</td>
            <td>window alpha (default 0.5)</td>
        </tr>
        <tr>
            <td>-h</td>
            <td>float</td>
            <td>window height (default 1080)</td>
        </tr>
        <tr>
            <td>-i</td>
            <td>string</td>
            <td>background image</td>
        </tr>
        <tr>
            <td>-u</td>
            <td>string</td>
            <td>login user</td>
        </tr>
        <tr>
            <td>-w</td>
            <td>float</td>
            <td>window width (default 1920)</td>
        </tr>
    </tbody>
</table>

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
