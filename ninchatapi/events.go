package ninchatapi

import (
	"github.com/ninchat/ninchat-go"
)

// Event interface is implemented by all event structs.
type Event interface {
	// MergeFrom fills in the parameters specified by the clientEvent.  An
	// UnexpectedEventError may be returned.
	MergeFrom(clientEvent *ninchat.Event) error
}

// NewSessionCreated.
func NewSessionCreated(clientEvent *ninchat.Event) (event *SessionCreated, err error) {
	if clientEvent != nil {
		event = new(SessionCreated)
		err = event.MergeFrom(clientEvent)
	}
	return
}

// MemberJoined is a union of target-specific events.
type MemberJoined struct {
	Channel *ChannelMemberJoined
	Queue   *QueueMemberJoined
	Realm   *RealmMemberJoined
}

// MergeFrom fills in the parameters specified by the clientEvent.  An
// UnexpectedEventError is returned if its type is not "channel_member_joined",
// "queue_member_joined" or "realm_member_joined".
func (union *MemberJoined) MergeFrom(clientEvent *ninchat.Event) (err error) {
	switch clientEvent.String() {
	case "channel_member_joined":
		return union.Channel.MergeFrom(clientEvent)

	case "queue_member_joined":
		return union.Queue.MergeFrom(clientEvent)

	case "realm_member_joined":
		return union.Realm.MergeFrom(clientEvent)

	default:
		return newError(clientEvent)
	}
}

// MemberParted is a union of target-specific events.
type MemberParted struct {
	Channel *ChannelMemberParted
	Queue   *QueueMemberParted
	Realm   *RealmMemberParted
}

// MergeFrom fills in the parameters specified by the clientEvent.  An
// UnexpectedEventError is returned if its type is not "channel_member_parted",
// "queue_member_parted" or "realm_member_parted".
func (union *MemberParted) MergeFrom(clientEvent *ninchat.Event) (err error) {
	switch clientEvent.String() {
	case "channel_member_parted":
		return union.Channel.MergeFrom(clientEvent)

	case "queue_member_parted":
		return union.Queue.MergeFrom(clientEvent)

	case "realm_member_parted":
		return union.Realm.MergeFrom(clientEvent)

	default:
		return newError(clientEvent)
	}
}

// MemberUpdated is a union of target-specific events.
type MemberUpdated struct {
	Channel *ChannelMemberUpdated
	Realm   *RealmMemberUpdated
}

// MergeFrom fills in the parameters specified by the clientEvent.  An
// UnexpectedEventError is returned if its type is not "channel_member_updated"
// or "realm_member_updated".
func (union *MemberUpdated) MergeFrom(clientEvent *ninchat.Event) (err error) {
	switch clientEvent.String() {
	case "channel_member_updated":
		return union.Channel.MergeFrom(clientEvent)

	case "realm_member_updated":
		return union.Realm.MergeFrom(clientEvent)

	default:
		return newError(clientEvent)
	}
}

// ChannelMembers parameter type.
type ChannelMembers map[string]*ChannelMemberEntry

func (target *ChannelMembers) MergeFrom(source map[string]interface{}) {
	// TODO
}

type ChannelMemberEntry struct {
	MemberAttrs map[string]*ChannelMemberAttrs `json:"member_attrs"`
	UserAttrs   *UserAttrs                     `json:"user_attrs"`
}

// Channels parameter type.
type Channels map[string]*ChannelResult

func (target *Channels) MergeFrom(source map[string]interface{}) {
	// TODO
}

type ChannelResult struct {
	ChannelAttrs *ChannelAttrs `json:"channel_attrs"`
	RealmId      *string       `json:"realm_id,omitempty"`
	Weight       float64       `json:"weight"`
}

// DialogueMembers parameter type.
type DialogueMembers map[string]*DialogueMemberAttrs

func (target *DialogueMembers) MergeFrom(source map[string]interface{}) {
	// TODO
}

// MasterKeys parameter type.
type MasterKeys map[string]struct{}

func (target *MasterKeys) MergeFrom(source map[string]interface{}) {
	// TODO
}

// QueueMembers parameter type.
type QueueMembers map[string]*QueueMemberEntry

func (target *QueueMembers) MergeFrom(source map[string]interface{}) {
	// TODO
}

type QueueMemberEntry struct {
	MemberAttrs map[string]struct{} `json:"member_attrs"`
	UserAttrs   *UserAttrs          `json:"user_attrs"`
}

// RealmMembers parameter type.
type RealmMembers map[string]*RealmMemberEntry

func (target *RealmMembers) MergeFrom(source map[string]interface{}) {
	// TODO
}

type RealmMemberEntry struct {
	MemberAttrs map[string]*RealmMemberAttrs `json:"member_attrs"`
	UserAttrs   *UserAttrs                   `json:"user_attrs"`
}

// RealmQueues parameter type.
type RealmQueues map[string]*RealmQueueEntry

func (target *RealmQueues) MergeFrom(source map[string]interface{}) {
	// TODO
}

type RealmQueueEntry struct {
	QueueAttrs    *QueueAttrs `json:"queue_attrs"`
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

func (target *UserAccount) MergeFrom(source map[string]interface{}) {
	// TODO
}

type UserAccountMembers struct {
	Quota int `json:"quota"`
}

type UserAccountObjects struct {
	Available int `json:"available"`
	Quota     int `json:"quota"`
}

type UserAccountSubscription struct {
	Active       bool                `json:"active,omitempty"`
	Channels     *UserAccountObjects `json:"channels,omitempty"`
	Expiration   *int                `json:"expiration,omitempty"`
	Plan         string              `json:"plan"`
	QueueMembers *UserAccountMembers `json:"queue_members,omitempty"`
	Queues       *UserAccountObjects `json:"queues,omitempty"`
	Realms       *UserAccountObjects `json:"realms,omitempty"`
	Renewal      *int                `json:"renewal,omitempty"`
}

// UserChannels parameter type.
type UserChannels map[string]*UserChannelEntry

func (target *UserChannels) MergeFrom(source map[string]interface{}) {
	// TODO
}

type UserChannelEntry struct {
	ChannelAttrs  *ChannelAttrs `json:"channel_attrs"`
	ChannelStatus *string       `json:"channel_status,omitempty"`
	RealmId       *string       `json:"realm_id,omitempty"`
}

// UserDialogues parameter type.
type UserDialogues map[string]*UserDialogueEntry

func (target *UserDialogues) MergeFrom(source map[string]interface{}) {
	// TODO
}

type UserDialogueEntry struct {
	AudienceMetadata map[string]interface{} `json:"audience_metadata,omitempty"`
	DialogueMembers  *DialogueMembers       `json:"dialogue_members,omitempty"`
	DialogueStatus   *string                `json:"dialogue_status,omitempty"`
}

// UserIdentities parameter type.
type UserIdentities map[string]map[string]*IdentityAttrs

func (target *UserIdentities) MergeFrom(source map[string]interface{}) {
	// TODO
}

// UserQueues parameter type.
type UserQueues map[string]*UserQueueEntry

func (target *UserQueues) MergeFrom(source map[string]interface{}) {
	// TODO
}

type UserQueueEntry struct {
	QueueAttrs *QueueAttrs `json:"queue_attrs"`
	RealmId    string      `json:"realm_id"`
}

// UserRealms parameter type.
type UserRealms map[string]*RealmAttrs

func (target *UserRealms) MergeFrom(source map[string]interface{}) {
	// TODO
}

// UserRealmsMember parameter type.
type UserRealmsMember map[string]*RealmMemberAttrs

func (target *UserRealmsMember) MergeFrom(source map[string]interface{}) {
	// TODO
}

// Users parameter type.
type Users map[string]*UserResult

func (target *Users) MergeFrom(source map[string]interface{}) {
	// TODO
}

type UserResult struct {
	UserAttrs *UserAttrs `json:"user_attrs"`
	Weight    float64    `json:"weight"`
}
