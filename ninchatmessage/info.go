package ninchatmessage

import (
	"github.com/ninchat/ninchat-go"
)

const (
	UserInfoType    = "ninchat.com/info/user"
	ChannelInfoType = "ninchat.com/info/channel"
	JoinInfoType    = "ninchat.com/info/join"
	PartInfoType    = "ninchat.com/info/part"
	MemberInfoType  = "ninchat.com/info/member"
	AccessInfoType  = "ninchat.com/info/access"

	// Pattern for all info messages (including ones added in future Ninchat
	// API updates).  For use with the "message_types" parameter passed to
	// ninchat.Session.SetParams().
	InfoTypes = "ninchat.com/info/*"
)

// UserInfo represents https://ninchat.com/info/user messages.
type UserInfo struct {
	UserId      string  `json:"user_id"`
	UserName    *string `json:"user_name,omitempty"`
	UserNameOld *string `json:"user_name_old,omitempty"`
	UserDeleted bool    `json:"user_deleted,omitempty"`
}

func (*UserInfo) MessageType() string {
	return UserInfoType
}

func (m *UserInfo) Unmarshal(payload []ninchat.Frame) (err error) {
	x, err := unmarshal(payload)
	if err == nil {
		m.Init(x)
	}
	return
}

func (m *UserInfo) Init(payload interface{}) {
	if obj, ok := payload.(map[string]interface{}); ok {
		if x := obj["user_id"]; x != nil {
			m.UserId, _ = x.(string)
		}

		if x := obj["user_name"]; x != nil {
			if y, ok := x.(string); ok {
				m.UserName = &y
			}
		}

		if x := obj["user_name_old"]; x != nil {
			if y, ok := x.(string); ok {
				m.UserNameOld = &y
			}
		}

		if obj["user_deleted"] != nil {
			m.UserDeleted = true
		}
	}
}

// ChannelInfo represents https://ninchat.com/info/channel messages.
type ChannelInfo struct {
	ChannelAttrsOld map[string]interface{} `json:"channel_attrs_old"`
	ChannelAttrsNew map[string]interface{} `json:"channel_attrs_new"`
}

func (*ChannelInfo) MessageType() string {
	return ChannelInfoType
}

func (m *ChannelInfo) Unmarshal(payload []ninchat.Frame) (err error) {
	x, err := unmarshal(payload)
	if err == nil {
		m.Init(x)
	}
	return
}

func (m *ChannelInfo) Init(payload interface{}) {
	if obj, ok := payload.(map[string]interface{}); ok {
		if x := obj["channel_attrs_old"]; x != nil {
			m.ChannelAttrsOld, _ = x.(map[string]interface{})
		}

		if x := obj["channel_attrs_new"]; x != nil {
			m.ChannelAttrsNew, _ = x.(map[string]interface{})
		}
	}
}

// JoinInfo represents https://ninchat.com/info/join messages.
type JoinInfo struct {
	UserId         string  `json:"user_id"`
	UserName       *string `json:"user_name,omitempty"`
	MemberSilenced bool    `json:"member_silenced,omitempty"`
}

func (*JoinInfo) MessageType() string {
	return JoinInfoType
}

func (m *JoinInfo) Unmarshal(payload []ninchat.Frame) (err error) {
	x, err := unmarshal(payload)
	if err == nil {
		m.Init(x)
	}
	return
}

func (m *JoinInfo) Init(payload interface{}) {
	if obj, ok := payload.(map[string]interface{}); ok {
		if x := obj["user_id"]; x != nil {
			m.UserId, _ = x.(string)
		}

		if x := obj["user_name"]; x != nil {
			if y, ok := x.(string); ok {
				m.UserName = &y
			}
		}

		if obj["member_silenced"] != nil {
			m.MemberSilenced = true
		}
	}
}

// PartInfo represents https://ninchat.com/info/part messages.
type PartInfo struct {
	UserId   string  `json:"user_id"`
	UserName *string `json:"user_name,omitempty"`
}

func (*PartInfo) MessageType() string {
	return PartInfoType
}

func (m *PartInfo) Unmarshal(payload []ninchat.Frame) (err error) {
	x, err := unmarshal(payload)
	if err == nil {
		m.Init(x)
	}
	return
}

func (m *PartInfo) Init(payload interface{}) {
	if obj, ok := payload.(map[string]interface{}); ok {
		if x := obj["user_id"]; x != nil {
			m.UserId, _ = x.(string)
		}

		if x := obj["user_name"]; x != nil {
			if y, ok := x.(string); ok {
				m.UserName = &y
			}
		}
	}
}

// MemberInfo represents https://ninchat.com/info/member messages.
type MemberInfo struct {
	UserId         string  `json:"user_id"`
	UserName       *string `json:"user_name,omitempty"`
	MemberSilenced bool    `json:"member_silenced,omitempty"`
}

func (*MemberInfo) MessageType() string {
	return MemberInfoType
}

func (m *MemberInfo) Unmarshal(payload []ninchat.Frame) (err error) {
	x, err := unmarshal(payload)
	if err == nil {
		m.Init(x)
	}
	return
}

func (m *MemberInfo) Init(payload interface{}) {
	if obj, ok := payload.(map[string]interface{}); ok {
		if x := obj["user_id"]; x != nil {
			m.UserId, _ = x.(string)
		}

		if x := obj["user_name"]; x != nil {
			if y, ok := x.(string); ok {
				m.UserName = &y
			}
		}

		if obj["member_silenced"] != nil {
			m.MemberSilenced = true
		}
	}
}

// AccessInfo represents https://ninchat.com/info/access messages.
type AccessInfo struct {
	UserId       string                 `json:"user_id"`
	AccessKey    string                 `json:"access_key"`
	ChannelId    string                 `json:"channel_id"`
	ChannelAttrs map[string]interface{} `json:"channel_attrs"`
	RealmId      *string                `json:"realm_id,omitempty"`
	RealmAttrs   map[string]interface{} `json:"realm_attrs,omitempty"`
	RealmMember  bool                   `json:"realm_member,omitempty"`
}

func (*AccessInfo) MessageType() string {
	return AccessInfoType
}

func (m *AccessInfo) Unmarshal(payload []ninchat.Frame) (err error) {
	x, err := unmarshal(payload)
	if err == nil {
		m.Init(x)
	}
	return
}

func (m *AccessInfo) Init(payload interface{}) {
	if obj, ok := payload.(map[string]interface{}); ok {
		if x := obj["user_id"]; x != nil {
			m.UserId, _ = x.(string)
		}

		if x := obj["access_key"]; x != nil {
			m.AccessKey, _ = x.(string)
		}

		if x := obj["channel_id"]; x != nil {
			m.ChannelId, _ = x.(string)
		}

		if x := obj["channel_attrs"]; x != nil {
			m.ChannelAttrs, _ = x.(map[string]interface{})
		}

		if x := obj["realm_id"]; x != nil {
			if y, ok := x.(string); ok {
				m.RealmId = &y
			}
		}

		if x := obj["realm_attrs"]; x != nil {
			m.RealmAttrs, _ = x.(map[string]interface{})
		}

		if obj["realm_member"] != nil {
			m.RealmMember = true
		}
	}
}
