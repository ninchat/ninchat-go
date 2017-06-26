// +build !js

package ninchat

import (
	"time"
)

type duration time.Duration

const (
	second      = duration(time.Second)
	millisecond = duration(time.Millisecond)
)

type timeTime time.Time

func timeNow() timeTime {
	return timeTime(time.Now())
}

func timeAdd(t timeTime, d duration) timeTime {
	return timeTime(time.Time(t).Add(time.Duration(d)))
}

func timeSub(t1, t2 timeTime) duration {
	return duration(time.Time(t1).Sub(time.Time(t2)))
}

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
