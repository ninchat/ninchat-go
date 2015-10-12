package ninchatapi

import (
	"github.com/ninchat/ninchat-go"
)

// Event interface is implemented by all event structs.
type Event interface {
	// MergeFrom fills in the parameters specified by the clientEvent.  An
	// UnexpectedEventError may be returned.
	Init(clientEvent *ninchat.Event) error
}

// MemberJoined is a union of target-specific events.
type MemberJoined struct {
	Channel *ChannelMemberJoined
	Queue   *QueueMemberJoined
	Realm   *RealmMemberJoined
}

// Init fills in the parameters specified by the clientEvent.  An
// UnexpectedEventError is returned if its type is not "channel_member_joined",
// "queue_member_joined" or "realm_member_joined".
func (union *MemberJoined) Init(clientEvent *ninchat.Event) (err error) {
	switch clientEvent.String() {
	case "channel_member_joined":
		return union.Channel.Init(clientEvent)

	case "queue_member_joined":
		return union.Queue.Init(clientEvent)

	case "realm_member_joined":
		return union.Realm.Init(clientEvent)

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

// Init fills in the parameters specified by the clientEvent.  An
// UnexpectedEventError is returned if its type is not "channel_member_parted",
// "queue_member_parted" or "realm_member_parted".
func (union *MemberParted) Init(clientEvent *ninchat.Event) (err error) {
	switch clientEvent.String() {
	case "channel_member_parted":
		return union.Channel.Init(clientEvent)

	case "queue_member_parted":
		return union.Queue.Init(clientEvent)

	case "realm_member_parted":
		return union.Realm.Init(clientEvent)

	default:
		return newError(clientEvent)
	}
}

// MemberUpdated is a union of target-specific events.
type MemberUpdated struct {
	Channel *ChannelMemberUpdated
	Realm   *RealmMemberUpdated
}

// Init fills in the parameters specified by the clientEvent.  An
// UnexpectedEventError is returned if its type is not "channel_member_updated"
// or "realm_member_updated".
func (union *MemberUpdated) Init(clientEvent *ninchat.Event) (err error) {
	switch clientEvent.String() {
	case "channel_member_updated":
		return union.Channel.Init(clientEvent)

	case "realm_member_updated":
		return union.Realm.Init(clientEvent)

	default:
		return newError(clientEvent)
	}
}
