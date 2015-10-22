package ninchatmessage

import (
	"github.com/ninchat/ninchat-go"
)

// Factories contains default constructors for all known message types.
var Factories = map[string]func() UnaryContent{
	AccessInfoType:  func() UnaryContent { return new(AccessInfo) },
	ChannelInfoType: func() UnaryContent { return new(ChannelInfo) },
	JoinInfoType:    func() UnaryContent { return new(JoinInfo) },
	LinkType:        func() UnaryContent { return new(Link) },
	MemberInfoType:  func() UnaryContent { return new(MemberInfo) },
	NoticeType:      func() UnaryContent { return new(Notice) },
	PartInfoType:    func() UnaryContent { return new(PartInfo) },
	TextType:        func() UnaryContent { return new(Text) },
	UserInfoType:    func() UnaryContent { return new(UserInfo) },
}

// NewContent creates a message object from a payload.  nil is returned if the
// type is unknown.
func NewContent(messageType string, payload []ninchat.Frame) (c Content, err error) {
	if f := Factories[messageType]; f != nil {
		c = f()
		err = c.Unmarshal(payload)
	}
	return
}

// NewUnaryContent creates a message object from an already parsed JSON value.
// nil is returned if the type is unknown.
func NewUnaryContent(messageType string, payload interface{}) (uc UnaryContent) {
	if f := Factories[messageType]; f != nil {
		uc = f()
		uc.Init(payload)
	}
	return
}
