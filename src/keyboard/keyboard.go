package keyboard

import (
	"log"

	// "github.com/linuxdeepin/dde-daemon/keybinding/shortcuts"
	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/keybind"

	"github.com/linuxdeepin/go-x11-client/util/keysyms"
	"github.com/linuxdeepin/go-x11-client/util/mousebind"
)

func SelectKeystroke(inputChan chan x.GenericEvent) error {
	conn, err := x.NewConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = grabKbdAndMouse(conn)
	if err != nil {
		log.Print("failed to grab keyboard and mouse:", err)
		return err
	}
	defer ungrabKbdAndMouse(conn)

	eventChan := make(chan x.GenericEvent, 500)
	conn.AddEventChan(eventChan)

	// keySymbols := keysyms.NewKeySymbols(conn)
	// var grabScreenKeystroke *shortcuts.Keystroke

	log.Printf("eventChan: %#v", eventChan)
	for evt := range eventChan {
		// log.Printf("event: %#v", evt)
		inputChan <- evt
		switch evt.GetEventCode() {
		case x.KeyPressEventCode:
			// event, _ := x.NewKeyPressEvent(evt)
			// log.Printf("KeyPressEventCode: %#v\n", event)
			sendEvent(conn, x.KeyPressEventCode, evt)
			// mods := shortcuts.GetConcernedModifiers(event.State)
			// log.Print("event mods:", shortcuts.Modifiers(event.State))
			// key := shortcuts.Key{
			// 	Mods: mods,
			// 	Code: shortcuts.Keycode(event.Detail),
			// }
			// log.Print("event key:", key)
			// ks := key.ToKeystroke(keySymbols)
			// emitSignalKeyEvent(true, ks.String())
			// if ks.IsGood() {
			// 	log.Print("good keystroke", ks)
			// 	grabScreenKeystroke = ks
			// } else {
			// 	log.Print("bad keystroke", ks)
			// 	grabScreenKeystroke = nil
			// }
		case x.KeyReleaseEventCode:
			// event, _ := x.NewKeyReleaseEvent(evt)
			// log.Printf("KeyReleaseEventCode: %#v\n", event)
			sendEvent(conn, x.KeyReleaseEventCode, evt)
			// if grabScreenKeystroke != nil {
			// 	emitSignalKeyEvent(false, grabScreenKeystroke.String())
			// 	grabScreenKeystroke = nil
			// } else {
			// 	emitSignalKeyEvent(false, "")
			// }

		case x.ButtonPressEventCode:
			log.Println("ButtonPressEvent")
			emitSignalKeyEvent(true, "")
		case x.ButtonReleaseEventCode:
			log.Println("ButtonReleaseEvent")
			emitSignalKeyEvent(false, "")
		}
	}

	log.Print("end selectKeystroke")
	return nil
}

func emitSignalKeyEvent(pressed bool, keystroke string) {
	log.Printf("pass key %s\n", keystroke)
	// m.service.Emit(m, "KeyEvent", pressed, keystroke)
}

func grabKbdAndMouse(conn *x.Conn) error {
	rootWin := conn.GetDefaultScreen().Root
	err := keybind.GrabKeyboard(conn, rootWin)
	if err != nil {
		return err
	}

	// Ignore mouse grab error
	const pointerEventMask = x.EventMaskButtonRelease | x.EventMaskButtonPress
	err = mousebind.GrabPointer(conn, rootWin, pointerEventMask, x.None, x.None)
	if err != nil {
		keybind.UngrabKeyboard(conn)
		return err
	}
	return nil
}

func ungrabKbdAndMouse(conn *x.Conn) {
	keybind.UngrabKeyboard(conn)
	mousebind.UngrabPointer(conn)
}

type CapsLockState uint

const (
	CapsLockOff CapsLockState = iota
	CapsLockOn
	CapsLockUnknown
)

func QueryCapsLockState(conn *x.Conn) (CapsLockState, error) {
	rootWin := conn.GetDefaultScreen().Root
	queryPointerReply, err := x.QueryPointer(conn, rootWin).Reply(conn)
	if err != nil {
		return CapsLockUnknown, err
	}

	// fmt.Printf("query pointer reply %#v", queryPointerReply)
	on := queryPointerReply.Mask&x.ModMaskLock != 0
	if on {
		return CapsLockOn, nil
	} else {
		return CapsLockOff, nil
	}
}

type ShiftState uint

const (
	ShiftOff ShiftState = iota
	ShiftOn
	ShiftUnknown
)

func QueryShiftState(conn *x.Conn) (ShiftState, error) {
	rootWin := conn.GetDefaultScreen().Root
	queryPointerReply, err := x.QueryPointer(conn, rootWin).Reply(conn)
	if err != nil {
		return ShiftUnknown, err
	}

	// fmt.Printf("query pointer reply %#v", queryPointerReply)
	on := queryPointerReply.Mask&x.ModMaskShift != 0
	if on {
		return ShiftOn, nil
	} else {
		return ShiftOff, nil
	}
}

func KeyCodeToString(keyCode x.Keycode, modifier uint16, conn *x.Conn) (string, bool) {
	keySymbols := keysyms.NewKeySymbols(conn)
	str, ok := keySymbols.LookupString(keyCode, modifier)
	if !ok {
		return str, ok
	}

	switch str {
	case "space":
		return " ", ok
	}

	return str, ok
}

func sendEvent(conn *x.Conn, eventCode uint32, event x.GenericEvent) error {
	log.Printf("pass key %#v (%v)\n", event, eventCode)
	rootWin := conn.GetDefaultScreen().Root
	x.SendEvent(conn, true, rootWin, eventCode, event)

	return nil
}
