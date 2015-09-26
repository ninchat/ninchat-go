package ninchatmessage

import (
	"encoding/json"
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

func (m *UserInfo) Marshal() (payload [][]byte, err error) {
	return marshalJSON(m)
}

func (m *UserInfo) Unmarshal(payload [][]byte) error {
	return json.Unmarshal(payload[0], m)
}

// ChannelInfo represents https://ninchat.com/info/channel messages.
type ChannelInfo struct {
	ChannelAttrsOld map[string]interface{} `json:"channel_attrs_old"`
	ChannelAttrsNew map[string]interface{} `json:"channel_attrs_new"`
}

func (*ChannelInfo) MessageType() string {
	return ChannelInfoType
}

func (m *ChannelInfo) Marshal() (payload [][]byte, err error) {
	return marshalJSON(m)
}

func (m *ChannelInfo) Unmarshal(payload [][]byte) error {
	return json.Unmarshal(payload[0], m)
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

func (m *JoinInfo) Marshal() (payload [][]byte, err error) {
	return marshalJSON(m)
}

func (m *JoinInfo) Unmarshal(payload [][]byte) error {
	return json.Unmarshal(payload[0], m)
}

// PartInfo represents https://ninchat.com/info/part messages.
type PartInfo struct {
	UserId   string  `json:"user_id"`
	UserName *string `json:"user_name,omitempty"`
}

func (*PartInfo) MessageType() string {
	return PartInfoType
}

func (m *PartInfo) Marshal() (payload [][]byte, err error) {
	return marshalJSON(m)
}

func (m *PartInfo) Unmarshal(payload [][]byte) error {
	return json.Unmarshal(payload[0], m)
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

func (m *MemberInfo) Marshal() (payload [][]byte, err error) {
	return marshalJSON(m)
}

func (m *MemberInfo) Unmarshal(payload [][]byte) error {
	return json.Unmarshal(payload[0], m)
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

func (m *AccessInfo) Marshal() (payload [][]byte, err error) {
	return marshalJSON(m)
}

func (m *AccessInfo) Unmarshal(payload [][]byte) error {
	return json.Unmarshal(payload[0], m)
}
