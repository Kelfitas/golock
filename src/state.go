package main

import "time"

type Event uint32

const (
	NoEvent Event = iota
	EmptyPasswordEvent
	CapsChangedEvent
	AuthSuccessEvent
	AuthCheckEvent
	WrongPasswordEvent
	BackSpaceEvent
	KeyPressEvent
)

type State struct {
	IsCapsLockOn     bool
	EventChan        chan Event
	LastEvent        Event
	ShouldDraw       bool
	LastQueuedRedraw time.Time
	LastStart        int32
	PasswordLength   int32
}

var state *State

func init() {
	eventChan := make(chan Event, 100)
	state = &State{EventChan: eventChan}
}
