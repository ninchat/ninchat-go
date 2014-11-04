package main

import (
	"github.com/gopherjs/gopherjs/js"
)

// Time
type Time int64

// Duration
type Duration int64

const (
	Millisecond Duration = 1
	Second               = Millisecond * 1000
	Minute               = Second * 60
)

// Now
func Now() Time {
	return Time(js.Global.Get("Date").New().Call("getTime").Int64())
}

// Timer
type Timer struct {
	C chan bool

	id js.Object
}

func NewTimer(timeout Duration) (timer *Timer) {
	timer = &Timer{
		C: make(chan bool),
	}

	if timeout >= 0 {
		timer.Reset(timeout)
	}

	return
}

func (timer *Timer) Active() bool {
	return timer.id != nil
}

func (timer *Timer) Reset(timeout Duration) {
	timer.Stop()

	timer.id = SetTimeout(func() {
		timer.id = nil

		go func() {
			timer.C <- true
		}()
	}, timeout)
}

func (timer *Timer) Stop() {
	if timer.id != nil {
		ClearTimeout(timer.id)
		timer.id = nil
	}
}

// Sleep
func Sleep(delay Duration) {
	<-NewTimer(delay).C
}
