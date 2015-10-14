package ninchatmessage_test

import (
	"testing"

	"."
)

func TestOriginalContent(t *testing.T) {
	_ = []ninchatmessage.OriginalContent{
		new(ninchatmessage.Link),
		new(ninchatmessage.Notice),
		new(ninchatmessage.Text),
	}
}
