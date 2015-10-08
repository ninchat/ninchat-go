package ninchatapi

import (
	ninchat ".."
)

type eventInit interface {
	init(*ninchat.Event) error
}

// MemberJoined is a union of target-specific events.
type MemberJoined struct {
	Channel *ChannelMemberJoined
	Queue   *QueueMemberJoined
	Realm   *RealmMemberJoined
}

func (union *MemberJoined) init(clientEvent *ninchat.Event) (err error) {
	switch clientEvent.String() {
	case "channel_member_joined":
		return union.Channel.init(clientEvent)

	case "queue_member_joined":
		return union.Queue.init(clientEvent)

	case "realm_member_joined":
		return union.Realm.init(clientEvent)

	default:
		return &EventError{clientEvent}
	}
}

// MemberParted is a union of target-specific events.
type MemberParted struct {
	Channel *ChannelMemberParted
	Queue   *QueueMemberParted
	Realm   *RealmMemberParted
}

func (union *MemberParted) init(clientEvent *ninchat.Event) (err error) {
	switch clientEvent.String() {
	case "channel_member_parted":
		return union.Channel.init(clientEvent)

	case "queue_member_parted":
		return union.Queue.init(clientEvent)

	case "realm_member_parted":
		return union.Realm.init(clientEvent)

	default:
		return &EventError{clientEvent}
	}
}

// MemberUpdated is a union of target-specific events.
type MemberUpdated struct {
	Channel *ChannelMemberUpdated
	Realm   *RealmMemberUpdated
}

func (union *MemberUpdated) init(clientEvent *ninchat.Event) (err error) {
	switch clientEvent.String() {
	case "channel_member_updated":
		return union.Channel.init(clientEvent)

	case "realm_member_updated":
		return union.Realm.init(clientEvent)

	default:
		return &EventError{clientEvent}
	}
}

// ChannelMembers parameter type.
type ChannelMembers map[string]*ChannelMemberEntry

func (target *ChannelMembers) init(source map[string]interface{}) {
	// TODO
}

type ChannelMemberEntry struct {
	MemberAttrs map[string]*ChannelMemberAttrs `json:"member_attrs,omitempty"`
	UserAttrs   *UserAttrs                     `json:"user_attrs,omitempty"`
}

// Channels parameter type.
type Channels map[string]*ChannelResult

func (target *Channels) init(source map[string]interface{}) {
	// TODO
}

type ChannelResult struct {
	ChannelAttrs *ChannelAttrs `json:"channel_attrs,omitempty"`
	RealmId      *string       `json:"realm_id,omitempty"`
	Weight       *float64      `json:"weight,omitempty"`
}

// DialogueMembers parameter type.
type DialogueMembers map[string]*DialogueMemberAttrs

func (target *DialogueMembers) init(source map[string]interface{}) {
	// TODO
}

// MasterKeys parameter type.
type MasterKeys map[string]struct{}

func (target *MasterKeys) init(source map[string]interface{}) {
	// TODO
}

// QueueMembers parameter type.
type QueueMembers map[string]*QueueMemberEntry

func (target *QueueMembers) init(source map[string]interface{}) {
	// TODO
}

type QueueMemberEntry struct {
	MemberAttrs map[string]struct{} `json:"member_attrs,omitempty"`
	UserAttrs   *UserAttrs          `json:"user_attrs,omitempty"`
}

// RealmMembers parameter type.
type RealmMembers map[string]*RealmMemberEntry

func (target *RealmMembers) init(source map[string]interface{}) {
	// TODO
}

type RealmMemberEntry struct {
	MemberAttrs map[string]*RealmMemberAttrs `json:"member_attrs,omitempty"`
	UserAttrs   *UserAttrs                   `json:"user_attrs,omitempty"`
}

// RealmQueues parameter type.
type RealmQueues map[string]*RealmQueueEntry

func (target *RealmQueues) init(source map[string]interface{}) {
	// TODO
}

type RealmQueueEntry struct {
	QueueAttrs    *QueueAttrs `json:"queue_attrs,omitempty"`
	QueuePosition *int        `json:"queue_position,omitempty"`
}

// UserAccount parameter type.
type UserAccount struct {
	Channels      *UserAccountObjects       `json:"channels,omitempty"`
	QueueMembers  *UserAccountMembers       `json:"queue_members,omitempty"`
	Queues        *UserAccountObjects       `json:"queues,omitempty"`
	Realms        *UserAccountObjects       `json:"realms,omitempty"`
	Subscriptions []UserAccountSubscription `json:"subscriptions,omitempty"`
}

func (target *UserAccount) init(source map[string]interface{}) {
	// TODO
}

type UserAccountMembers struct {
	Quota *int `json:"quota,omitempty"`
}

type UserAccountObjects struct {
	Available *int `json:"available,omitempty"`
	Quota     *int `json:"quota,omitempty"`
}

type UserAccountSubscription struct {
	Active       bool                `json:"active,omitempty"`
	Channels     *UserAccountObjects `json:"channels,omitempty"`
	Expiration   *int                `json:"expiration,omitempty"`
	Plan         *string             `json:"plan,omitempty"`
	QueueMembers *UserAccountMembers `json:"queue_members,omitempty"`
	Queues       *UserAccountObjects `json:"queues,omitempty"`
	Realms       *UserAccountObjects `json:"realms,omitempty"`
	Renewal      *int                `json:"renewal,omitempty"`
}

// UserChannels parameter type.
type UserChannels map[string]*UserChannelEntry

func (target *UserChannels) init(source map[string]interface{}) {
	// TODO
}

type UserChannelEntry struct {
	ChannelAttrs  *ChannelAttrs `json:"channel_attrs,omitempty"`
	ChannelStatus *string       `json:"channel_status,omitempty"`
	RealmId       *string       `json:"realm_id,omitempty"`
}

// UserDialogues parameter type.
type UserDialogues map[string]*UserDialogueEntry

func (target *UserDialogues) init(source map[string]interface{}) {
	// TODO
}

type UserDialogueEntry struct {
	AudienceMetadata map[string]interface{} `json:"audience_metadata,omitempty"`
	DialogueMembers  *DialogueMembers       `json:"dialogue_members,omitempty"`
	DialogueStatus   *string                `json:"dialogue_status,omitempty"`
}

// UserIdentities parameter type.
type UserIdentities map[string]map[string]*IdentityAttrs

func (target *UserIdentities) init(source map[string]interface{}) {
	// TODO
}

// UserQueues parameter type.
type UserQueues map[string]*UserQueueEntry

func (target *UserQueues) init(source map[string]interface{}) {
	// TODO
}

type UserQueueEntry struct {
	QueueAttrs *QueueAttrs `json:"queue_attrs,omitempty"`
	RealmId    *string     `json:"realm_id,omitempty"`
}

// UserRealms parameter type.
type UserRealms map[string]*RealmAttrs

func (target *UserRealms) init(source map[string]interface{}) {
	// TODO
}

// UserRealmsMember parameter type.
type UserRealmsMember map[string]*RealmMemberAttrs

func (target *UserRealmsMember) init(source map[string]interface{}) {
	// TODO
}

// Users parameter type.
type Users map[string]*UserResult

func (target *Users) init(source map[string]interface{}) {
	// TODO
}

type UserResult struct {
	UserAttrs *UserAttrs `json:"user_attrs,omitempty"`
	Weight    *float64   `json:"weight,omitempty"`
}
