package webhook

import (
	"encoding/json"
)

type MessageType string

const (
	MessageFile        MessageType = "ninchat.com/file"
	MessageChannelInfo MessageType = "ninchat.com/info/channel"
	MessageJoinInfo    MessageType = "ninchat.com/info/join"
	MessageMemberInfo  MessageType = "ninchat.com/info/member"
	MessagePartInfo    MessageType = "ninchat.com/info/part"
	MessageUserInfo    MessageType = "ninchat.com/info/user"
	MessageMetadata    MessageType = "ninchat.com/metadata"
	MessageNotice      MessageType = "ninchat.com/notice"
	MessageText        MessageType = "ninchat.com/text"
	MessageUIAction    MessageType = "ninchat.com/ui/action"
	MessageUICompose   MessageType = "ninchat.com/ui/compose"
)

type FilePayload struct {
	Text  *string `json:"text,omitempty"`
	Files []File  `json:"files"`
}

type File struct {
	FileID        string                     `json:"file_id"`
	FileAttrsJSON map[string]json.RawMessage `json:"file_attrs"`
}

func (m *Message) FilePayload() (p FilePayload, err error) {
	err = json.Unmarshal(m.PayloadJSON, &p)
	return
}

type ChannelInfoPayload struct {
	ChannelAttrsOldJSON map[string]json.RawMessage `json:"channel_attrs_old"`
	ChannelAttrsNewJSON map[string]json.RawMessage `json:"channel_attrs_new"`
}

func (m *Message) ChannelInfoPayload() (p ChannelInfoPayload, err error) {
	err = json.Unmarshal(m.PayloadJSON, &p)
	return
}

type JoinInfoPayload struct {
	UserID         string  `json:"user_id"`
	UserName       *string `json:"user_name,omitempty"`
	MemberSilenced bool    `json:"member_silenced,omitempty"`
}

func (m *Message) JoinInfoPayload() (p JoinInfoPayload, err error) {
	err = json.Unmarshal(m.PayloadJSON, &p)
	return
}

type MemberInfoPayload struct {
	UserID         string  `json:"user_id"`
	UserName       *string `json:"user_name,omitempty"`
	MemberSilenced bool    `json:"member_silenced,omitempty"`
}

func (m *Message) MemberInfoPayload() (p MemberInfoPayload, err error) {
	err = json.Unmarshal(m.PayloadJSON, &p)
	return
}

type PartInfoPayload struct {
	UserID   string  `json:"user_id"`
	UserName *string `json:"user_name,omitempty"`
	Cause    string  `json:"cause,omitempty"`
}

func (m *Message) PartInfoPayload() (p PartInfoPayload, err error) {
	err = json.Unmarshal(m.PayloadJSON, &p)
	return
}

type UserInfoPayload struct {
	UserID      string  `json:"user_id"`
	UserName    *string `json:"user_name,omitempty"`
	UserNameOld *string `json:"user_name_old,omitempty"`
	UserDeleted bool    `json:"user_deleted,omitempty"`
}

func (m *Message) UserInfoPayload() (p UserInfoPayload, err error) {
	err = json.Unmarshal(m.PayloadJSON, &p)
	return
}

type MetadataPayload struct {
	Data Metadata `json:"data"`
	Time float64  `json:"time,omitempty"`
}

func (m *Message) MetadataPayload() (p MetadataPayload, err error) {
	err = json.Unmarshal(m.PayloadJSON, &p)
	return
}

type NoticePayload struct {
	TextPayload
}

func (m *Message) NoticePayload() (p NoticePayload, err error) {
	err = json.Unmarshal(m.PayloadJSON, &p)
	return
}

type TextPayload struct {
	Text string `json:"text"`
}

func (m *Message) TextPayload() (p TextPayload, err error) {
	err = json.Unmarshal(m.PayloadJSON, &p)
	return
}

type UIActionPayload struct {
	Action string    `json:"action"`
	Target UICompose `json:"target"`
}

func (m *Message) UIActionPayload() (p UIActionPayload, err error) {
	err = json.Unmarshal(m.PayloadJSON, &p)
	return
}

type UIComposePayload []UICompose

type UICompose struct {
	Class   string `json:"class,omitempty"`
	Element string `json:"class"`
	ID      string `json:"id,omitempty"`
	Label   string `json:"label,omitempty"`
	Name    string `json:"name,omitempty"`
}

func (m *Message) UIComposePayload() (p UIComposePayload, err error) {
	err = json.Unmarshal(m.PayloadJSON, &p)
	return
}
