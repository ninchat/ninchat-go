package ninchatmodel

var numericStatuses = map[string]int{
	"hidden":    1,
	"visible":   2,
	"unread":    3,
	"highlight": 4,
}

func compareStatus(a, b string) int {
	return numericStatuses[a] - numericStatuses[b]
}
