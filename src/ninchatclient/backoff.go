package main

const (
	maxBackoffSlots = 1024
)

// Backoff
type Backoff struct {
	lastSlot int
}

func (b *Backoff) Success() {
	b.lastSlot = 0
}

func (b *Backoff) Failure(maxDelay Duration) (delay Duration) {
	if b.lastSlot > 0 {
		delay = Duration(Jitter(float64(maxDelay)*float64(b.lastSlot)/maxBackoffSlots, -0.5))
	}

	if b.lastSlot < maxBackoffSlots-1 {
		b.lastSlot = ((b.lastSlot + 1) << 1) - 1
	}

	return
}
