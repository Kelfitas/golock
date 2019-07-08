package auth

import (
	"errors"
	"fmt"

	"github.com/msteinert/pam"
)

// This example uses whatever default PAM service configuration is available
// on the system, and tries to authenticate any user. This should cause PAM
// to ask its conversation handler for a username and password, in sequence.
//
// This application will handle those requests by displaying the
// PAM-provided prompt and sending back the first line of stdin input
// it can read for each.
//
// Keep in mind that unless run as root (or setuid root), the only
// user's authentication that can succeed is that of the process owner.
func Check(user, password string) error {
	t, err := pam.StartFunc("golock", user, func(s pam.Style, msg string) (string, error) {
		fmt.Println("========================================")
		defer fmt.Println("========================================")
		switch s {
		case pam.PromptEchoOff:
			fmt.Printf("pam.PromptEchoOff: %#v\n", s)
			fmt.Println(msg)
			return password, nil
			// return speakeasy.Ask(msg)
		case pam.PromptEchoOn:
			fmt.Printf("pam.PromptEchoOn: %#v\n", s)
			fmt.Println(msg)
			return password, nil
		case pam.ErrorMsg:
			fmt.Printf("pam.ErrorMsg: %#v\n", s)
			fmt.Println(msg)
			return "", nil
		case pam.TextInfo:
			fmt.Printf("pam.TextInfo: %#v\n", s)
			fmt.Println(msg)
			return "", nil
		}

		fmt.Printf("Unrecognized message style: %#v\n", s)
		return "", errors.New("Unrecognized message style")
	})
	if err != nil {
		fmt.Printf("Start: %s", err.Error())
		return err
	}

	return t.Authenticate(0)
}
