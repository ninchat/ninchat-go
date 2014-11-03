package main

// Jitter randomizes a float by increasing or decreasing it.
func Jitter(x float64, scale float64) float64 {
	return x + x*scale*Random()
}

// JitterDuration randomizes a timeout by increasing or decreasing it.
func JitterDuration(d Duration, scale float64) Duration {
	return Duration(Jitter(float64(d), scale))
}

// JitterUint64 randomizes an integer by increasing or decreasing it.
func JitterUint64(n uint64, scale float64) uint64 {
	return uint64(Jitter(float64(n), scale))
}
