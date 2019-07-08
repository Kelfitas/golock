package main

type Event uint32

const (
	EmptyPasswordEvent Event = iota
	CapsChangedEvent
	WrongPasswordEvent
	KeyPressEvent
)

type State struct {
	IsCapsLockOn bool
	EventChan    chan Event
}

var state *State

func init() {
	eventChan := make(chan Event, 100)
	state = &State{EventChan: eventChan}
}
