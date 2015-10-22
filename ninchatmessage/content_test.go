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

func TestUnaryContent(t *testing.T) {
	_ = []ninchatmessage.UnaryContent{
		new(ninchatmessage.AccessInfo),
		new(ninchatmessage.ChannelInfo),
		new(ninchatmessage.JoinInfo),
		new(ninchatmessage.Link),
		new(ninchatmessage.MemberInfo),
		new(ninchatmessage.Notice),
		new(ninchatmessage.PartInfo),
		new(ninchatmessage.Text),
		new(ninchatmessage.UserInfo),
	}
}
