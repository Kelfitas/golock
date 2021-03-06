package main

import (
	"flag"
	"fmt"
	"golock/src/auth"
	"log"
	osUser "os/user"
	"time"

	"golock/src/keyboard"

	// "github.com/BurntSushi/xgb"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/keysyms"
)

var (
	user string
)

func init() {
	flag.StringVar(&user, "u", "", "login user")
}

func grabInputs(inputChan chan x.GenericEvent, errChan chan error) {
	go keyboard.SelectKeystroke(inputChan, errChan)
}

func isAllowedInPass(s string) bool {
	switch s {
	case " ":
		return true
	}

	sym, ok := keysyms.StringToKeysym(s)
	if !ok {
		return false
	}

	return !(keysyms.IsKeypadKey(sym) ||
		keysyms.IsPrivateKeypadKey(sym) ||
		keysyms.IsCursorKey(sym) ||
		keysyms.IsKeypadFuncationKey(sym) ||
		keysyms.IsFunctionKey(sym) ||
		keysyms.IsMiscFunctionKey(sym) ||
		keysyms.IsModifierKey(sym)) && len(s) == 1
}

func getPass(inputChan chan x.GenericEvent, errChan chan error) (string, error) {
	var pass string
	var isShiftOn bool
	var isCapsOn bool
	var isPasswordEmpty bool

	getMod := func(_isCapsOn, _isShiftOn bool) uint16 {
		var mods uint16
		if _isShiftOn {
			mods |= x.ModMaskShift
		}

		if _isCapsOn {
			mods |= x.ModMaskLock
		}

		return mods
	}

	conn, err := x.NewConn()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	for {
		select {
		case e := <-inputChan:
			switch e.GetEventCode() {
			case x.KeyPressEventCode:
				event, _ := x.NewKeyPressEvent(e)
				str, ok := keyboard.KeyCodeToString(event.Detail, getMod(isShiftOn, isCapsOn), conn)

				// Check if key is allowed in pass
				shouldAppend := isAllowedInPass(str)
				if shouldAppend {
					pass += str
					isPasswordEmpty = false
					state.EventChan <- KeyPressEvent
				} else {
					// check mods
					capsLockState, _ := keyboard.QueryCapsLockState(conn)
					if capsLockState == keyboard.CapsLockOn {
						isCapsOn = true
					} else {
						isCapsOn = false
					}
					if state.IsCapsLockOn != isCapsOn {
						state.IsCapsLockOn = isCapsOn
						state.EventChan <- CapsChangedEvent
					}

					shiftState, _ := keyboard.QueryShiftState(conn)
					if shiftState == keyboard.ShiftOn {
						isShiftOn = true
					} else {
						isShiftOn = false
					}

					sym, ok := keysyms.StringToKeysym(str)
					if ok {
						switch sym {
						case keysyms.XK_Return:
							return pass, nil
						case keysyms.XK_BackSpace:
							state.EventChan <- BackSpaceEvent
							if len(pass) > 0 {
								pass = pass[:len(pass)-1]
							}
						case keysyms.XK_Escape:
							pass = ""
						}
					}
				}

				state.PasswordLength = int32(len(pass))
				if len(pass) == 0 && !isPasswordEmpty {
					isPasswordEmpty = true
					state.EventChan <- EmptyPasswordEvent
					log.Println("Password empty!")
				}
				log.Printf("KeyPressEventCode: %s | %t (%#v)\n", str, ok, event.Detail)
			case x.KeyReleaseEventCode:
				// event, _ := x.NewKeyReleaseEvent(e)
				// str, ok := keyboard.KeyCodeToString(event.Detail, 0, conn)
				// log.Printf("KeyReleaseEventCode: %s | %t (%#v)\n", str, ok, event.Detail)
			}
		case err := <-errChan:
			return pass, err
			// case <-time.After(3 * time.Second):
			// 	return pass, nil
		}
	}

	// return "", nil
}

func watchAuth(done chan bool) {
	if user == "" {
		u, err := osUser.Current()
		if err != nil {
			log.Fatal("Failed to get current user, try using -u flag. err:", err)
		}
		user = u.Username
	}

	inputChan := make(chan x.GenericEvent)
	errChan := make(chan error)
	grabInputs(inputChan, errChan)
	defer func() {
		log.Println("Done!")
		done <- true
	}()

	for {
		pass, err := getPass(inputChan, errChan)
		// fmt.Printf("pass: %#v | err: %#v\n", pass, err)
		// pass, err = speakeasy.Ask("Enter password: ")
		// if err != nil {
		// 	panic(err)
		// }

		state.EventChan <- AuthCheckEvent
		err = auth.Check(user, pass)
		fmt.Printf("%#v\n", err)
		if err == nil {
			state.EventChan <- AuthSuccessEvent
			fmt.Println("Auth success!")
			time.Sleep(1 * time.Second)
			return
		}
		fmt.Println("Auth Failed! Try again!")
		state.EventChan <- WrongPasswordEvent
		fmt.Println("after event sent")

		// time.Sleep(5 * time.Second)
		fmt.Println("after time.sleep")
	}
}
