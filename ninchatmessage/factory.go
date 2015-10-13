package ninchatmessage

import (
	"github.com/ninchat/ninchat-go"
)

// Factories contains default constructors for all known message types.
var Factories = map[string]func() Content{
	AccessInfoType:  func() Content { return new(AccessInfo) },
	ChannelInfoType: func() Content { return new(ChannelInfo) },
	JoinInfoType:    func() Content { return new(JoinInfo) },
	LinkType:        func() Content { return new(Link) },
	MemberInfoType:  func() Content { return new(MemberInfo) },
	NoticeType:      func() Content { return new(Notice) },
	PartInfoType:    func() Content { return new(PartInfo) },
	TextType:        func() Content { return new(Text) },
	UserInfoType:    func() Content { return new(UserInfo) },
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
