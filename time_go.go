// +build !js

package ninchat

import (
	"time"
)

type duration time.Duration

const (
	second = duration(time.Second)
)

type timer struct {
	C <-chan time.Time

	impl *time.Timer
}

func newTimer(timeout duration) *timer {
	impl := time.NewTimer(time.Duration(timeout))

	return &timer{
		C:    impl.C,
		impl: impl,
	}
}

func (timer *timer) Active() bool {
	return timer.impl != nil
}

func (timer *timer) Reset(timeout duration) {
	timer.impl.Reset(time.Duration(timeout))
}

func (timer *timer) Stop() {
	timer.impl.Stop()
	timer.impl = nil
}
