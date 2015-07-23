// +build js

package ninchat

import (
	"github.com/gopherjs/gopherjs/js"
)

type duration int64

const (
	second duration = 1000
)

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

func sleep(delay duration) {
	<-newTimer(delay).C
}
