package ninchatapi

// AppendStrings duplicates the source slice while unwrapping the elements.
func AppendStrings(target []string, source []interface{}) []string {
	if source != nil {
		if target == nil || cap(target) < len(target)+len(source) {
			t := make([]string, len(target), len(target)+len(source))
			copy(t, target)
			target = t
		}

		for _, x := range source {
			y, _ := x.(string)
			target = append(target, y)
		}
	}

	return target
}

func intPointer(x float64) *int {
	y := int(x)
	return &y
}
