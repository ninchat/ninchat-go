package ninchat

// jitterFloat64 randomizes a float by increasing or decreasing it.
func jitterFloat64(x float64, scale float64) float64 {
	return x + x*scale*randFloat64()
}

// jitterDuration randomizes a timeout by increasing or decreasing it.
func jitterDuration(d duration, scale float64) duration {
	return duration(jitterFloat64(float64(d), scale))
}

// jitterInt64 randomizes an integer by increasing or decreasing it.
func jitterInt64(n int64, scale float64) int64 {
	return int64(jitterFloat64(float64(n), scale))
}
