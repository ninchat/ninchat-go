// +build js

package ninchat

import (
	"github.com/gopherjs/gopherjs/js"
)

type duration int64

const (
	second      duration = 1000
	millisecond          = 1
)

type timeTime int64

func timeNow() timeTime {
	return timeTime(js.Global.Get("Date").New().Call("getTime").Int64())
}

func timeAdd(t timeTime, d duration) timeTime {
	return t + timeTime(d)
}

func timeSub(t1, t2 timeTime) duration {
	return duration(t1 - t2)
}

type timer struct {
	C chan struct{}

	id *js.Object
}

func newTimer(timeout duration) (t *timer) {
	t = &timer{
		C: make(chan struct{}),
	}

	if timeout >= 0 {
		t.Reset(timeout)
	}

	return
}

func (timer *timer) Active() bool {
	return timer.id != nil
}

func (timer *timer) Reset(timeout duration) {
	timer.Stop()

	timer.id = js.Global.Call("setTimeout", func() {
		timer.id = nil

		go func() {
			timer.C <- struct{}{}
		}()
	}, timeout)
}

func (timer *timer) Stop() {
	if timer.id != nil {
		js.Global.Call("clearTimeout", timer.id)
		timer.id = nil
	}
}
