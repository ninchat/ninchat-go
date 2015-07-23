package ninchat

const (
	maxBackoffSlots = 1024
)

type backoff struct {
	lastSlot int
}

func (b *backoff) success() {
	b.lastSlot = 0
}

func (b *backoff) failure(maxDelay duration) (delay duration) {
	if b.lastSlot > 0 {
		delay = duration(jitterFloat64(float64(maxDelay)*float64(b.lastSlot)/maxBackoffSlots, -0.5))
	}

	if b.lastSlot < maxBackoffSlots-1 {
		b.lastSlot = ((b.lastSlot + 1) << 1) - 1
	}

	return
}
