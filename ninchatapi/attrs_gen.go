package ninchatapi

// THIS FILE IS AUTO-GENERATED BY generate.py - DO NOT EDIT BY HAND!

// ChannelAttrs.  https://ninchat.com/api/v2#channel
type ChannelAttrs struct {
	Autohide                bool     `json:"autohide,omitempty"`
	Autosilence             bool     `json:"autosilence,omitempty"`
	BlacklistedMessageTypes []string `json:"blacklisted_message_types,omitempty"`
	Closed                  bool     `json:"closed,omitempty"`
	DisclosedSince          *int     `json:"disclosed_since,omitempty"`
	Followable              bool     `json:"followable,omitempty"`
	Name                    *string  `json:"name,omitempty"`
	OwnerId                 *string  `json:"owner_id,omitempty"`
	Private                 bool     `json:"private,omitempty"`
	Public                  bool     `json:"public,omitempty"`
	Ratelimit               *string  `json:"ratelimit,omitempty"`
	Suspended               bool     `json:"suspended,omitempty"`
	Topic                   *string  `json:"topic,omitempty"`
	Upload                  *string  `json:"upload,omitempty"`
	VerifiedJoin            bool     `json:"verified_join,omitempty"`
}

// NewChannelAttrs creates an object with the attributes specified by the source.
func NewChannelAttrs(source map[string]interface{}) (target *ChannelAttrs) {
	target = new(ChannelAttrs)
	target.Init(source)
	return
}

// Init fills in the attributes specified by the source
// (other fields are not touched).
func (target *ChannelAttrs) Init(source map[string]interface{}) {
	if x := source["autohide"]; x != nil {
		target.Autohide = true
	}

	if x := source["autosilence"]; x != nil {
		target.Autosilence = true
	}

	if x := source["blacklisted_message_types"]; x != nil {
		if y, ok := x.([]interface{}); ok {
			target.BlacklistedMessageTypes = AppendStrings(nil, y)
		}
	}

	if x := source["closed"]; x != nil {
		target.Closed = true
	}

	if x := source["disclosed_since"]; x != nil {
		if y, ok := x.(int); ok {
			target.DisclosedSince = &y
		}
	}

	if x := source["followable"]; x != nil {
		target.Followable = true
	}

	if x := source["name"]; x != nil {
		if y, ok := x.(string); ok {
			target.Name = &y
		}
	}

	if x := source["owner_id"]; x != nil {
		if y, ok := x.(string); ok {
			target.OwnerId = &y
		}
	}

	if x := source["private"]; x != nil {
		target.Private = true
	}

	if x := source["public"]; x != nil {
		target.Public = true
	}

	if x := source["ratelimit"]; x != nil {
		if y, ok := x.(string); ok {
			target.Ratelimit = &y
		}
	}

	if x := source["suspended"]; x != nil {
		target.Suspended = true
	}

	if x := source["topic"]; x != nil {
		if y, ok := x.(string); ok {
			target.Topic = &y
		}
	}

	if x := source["upload"]; x != nil {
		if y, ok := x.(string); ok {
			target.Upload = &y
		}
	}

	if x := source["verified_join"]; x != nil {
		target.VerifiedJoin = true
	}
}

// ChannelMemberAttrs.  https://ninchat.com/api/v2#channel-membership
type ChannelMemberAttrs struct {
	Autohide  bool `json:"autohide,omitempty"`
	Moderator bool `json:"moderator,omitempty"`
	Operator  bool `json:"operator,omitempty"`
	Silenced  bool `json:"silenced,omitempty"`
	Since     *int `json:"since,omitempty"`
}

// NewChannelMemberAttrs creates an object with the attributes specified by the source.
func NewChannelMemberAttrs(source map[string]interface{}) (target *ChannelMemberAttrs) {
	target = new(ChannelMemberAttrs)
	target.Init(source)
	return
}

// Init fills in the attributes specified by the source
// (other fields are not touched).
func (target *ChannelMemberAttrs) Init(source map[string]interface{}) {
	if x := source["autohide"]; x != nil {
		target.Autohide = true
	}

	if x := source["moderator"]; x != nil {
		target.Moderator = true
	}

	if x := source["operator"]; x != nil {
		target.Operator = true
	}

	if x := source["silenced"]; x != nil {
		target.Silenced = true
	}

	if x := source["since"]; x != nil {
		if y, ok := x.(int); ok {
			target.Since = &y
		}
	}
}

// DialogueMemberAttrs.  https://ninchat.com/api/v2#dialogue-membership
type DialogueMemberAttrs struct {
	AudienceEnded bool    `json:"audience_ended,omitempty"`
	QueueId       *string `json:"queue_id,omitempty"`
	Rating        *int    `json:"rating,omitempty"`
	Writing       bool    `json:"writing,omitempty"`
}

// NewDialogueMemberAttrs creates an object with the attributes specified by the source.
func NewDialogueMemberAttrs(source map[string]interface{}) (target *DialogueMemberAttrs) {
	target = new(DialogueMemberAttrs)
	target.Init(source)
	return
}

// Init fills in the attributes specified by the source
// (other fields are not touched).
func (target *DialogueMemberAttrs) Init(source map[string]interface{}) {
	if x := source["audience_ended"]; x != nil {
		target.AudienceEnded = true
	}

	if x := source["queue_id"]; x != nil {
		if y, ok := x.(string); ok {
			target.QueueId = &y
		}
	}

	if x := source["rating"]; x != nil {
		if y, ok := x.(int); ok {
			target.Rating = &y
		}
	}

	if x := source["writing"]; x != nil {
		target.Writing = true
	}
}

// IdentityAttrs.  https://ninchat.com/api/v2#identity
type IdentityAttrs struct {
	Auth     bool `json:"auth,omitempty"`
	Blocked  bool `json:"blocked,omitempty"`
	Pending  bool `json:"pending,omitempty"`
	Public   bool `json:"public,omitempty"`
	Rejected bool `json:"rejected,omitempty"`
}

// NewIdentityAttrs creates an object with the attributes specified by the source.
func NewIdentityAttrs(source map[string]interface{}) (target *IdentityAttrs) {
	target = new(IdentityAttrs)
	target.Init(source)
	return
}

// Init fills in the attributes specified by the source
// (other fields are not touched).
func (target *IdentityAttrs) Init(source map[string]interface{}) {
	if x := source["auth"]; x != nil {
		target.Auth = true
	}

	if x := source["blocked"]; x != nil {
		target.Blocked = true
	}

	if x := source["pending"]; x != nil {
		target.Pending = true
	}

	if x := source["public"]; x != nil {
		target.Public = true
	}

	if x := source["rejected"]; x != nil {
		target.Rejected = true
	}
}

// PuppetAttrs.  https://ninchat.com/api/v2#puppet
type PuppetAttrs struct {
	Name *string `json:"name,omitempty"`
}

// NewPuppetAttrs creates an object with the attributes specified by the source.
func NewPuppetAttrs(source map[string]interface{}) (target *PuppetAttrs) {
	target = new(PuppetAttrs)
	target.Init(source)
	return
}

// Init fills in the attributes specified by the source
// (other fields are not touched).
func (target *PuppetAttrs) Init(source map[string]interface{}) {
	if x := source["name"]; x != nil {
		if y, ok := x.(string); ok {
			target.Name = &y
		}
	}
}

// QueueAttrs.  https://ninchat.com/api/v2#queue
type QueueAttrs struct {
	Capacity  *int    `json:"capacity,omitempty"`
	Closed    bool    `json:"closed,omitempty"`
	Length    *int    `json:"length,omitempty"`
	Name      *string `json:"name,omitempty"`
	Suspended bool    `json:"suspended,omitempty"`
	Upload    *string `json:"upload,omitempty"`
}

// NewQueueAttrs creates an object with the attributes specified by the source.
func NewQueueAttrs(source map[string]interface{}) (target *QueueAttrs) {
	target = new(QueueAttrs)
	target.Init(source)
	return
}

// Init fills in the attributes specified by the source
// (other fields are not touched).
func (target *QueueAttrs) Init(source map[string]interface{}) {
	if x := source["capacity"]; x != nil {
		if y, ok := x.(int); ok {
			target.Capacity = &y
		}
	}

	if x := source["closed"]; x != nil {
		target.Closed = true
	}

	if x := source["length"]; x != nil {
		if y, ok := x.(int); ok {
			target.Length = &y
		}
	}

	if x := source["name"]; x != nil {
		if y, ok := x.(string); ok {
			target.Name = &y
		}
	}

	if x := source["suspended"]; x != nil {
		target.Suspended = true
	}

	if x := source["upload"]; x != nil {
		if y, ok := x.(string); ok {
			target.Upload = &y
		}
	}
}

// QueueMemberAttrs.
type QueueMemberAttrs struct {
}

// NewQueueMemberAttrs creates an object with the attributes specified by the source.
func NewQueueMemberAttrs(source map[string]interface{}) (target *QueueMemberAttrs) {
	target = new(QueueMemberAttrs)
	target.Init(source)
	return
}

// Init fills in the attributes specified by the source
// (other fields are not touched).
func (target *QueueMemberAttrs) Init(source map[string]interface{}) {
}

// RealmAttrs.  https://ninchat.com/api/v2#realm
type RealmAttrs struct {
	Name         *string                `json:"name,omitempty"`
	OwnerAccount *RealmOwnerAccountAttr `json:"owner_account,omitempty"`
	OwnerId      *string                `json:"owner_id,omitempty"`
	Suspended    bool                   `json:"suspended,omitempty"`
	Theme        *RealmThemeAttr        `json:"theme,omitempty"`
}

// NewRealmAttrs creates an object with the attributes specified by the source.
func NewRealmAttrs(source map[string]interface{}) (target *RealmAttrs) {
	target = new(RealmAttrs)
	target.Init(source)
	return
}

// Init fills in the attributes specified by the source
// (other fields are not touched).
func (target *RealmAttrs) Init(source map[string]interface{}) {
	if x := source["name"]; x != nil {
		if y, ok := x.(string); ok {
			target.Name = &y
		}
	}

	if x := source["owner_account"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.OwnerAccount = NewRealmOwnerAccountAttr(y)
		}
	}

	if x := source["owner_id"]; x != nil {
		if y, ok := x.(string); ok {
			target.OwnerId = &y
		}
	}

	if x := source["suspended"]; x != nil {
		target.Suspended = true
	}

	if x := source["theme"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.Theme = NewRealmThemeAttr(y)
		}
	}
}

// RealmMemberAttrs.  https://ninchat.com/api/v2#realm-membership
type RealmMemberAttrs struct {
	Operator bool `json:"operator,omitempty"`
}

// NewRealmMemberAttrs creates an object with the attributes specified by the source.
func NewRealmMemberAttrs(source map[string]interface{}) (target *RealmMemberAttrs) {
	target = new(RealmMemberAttrs)
	target.Init(source)
	return
}

// Init fills in the attributes specified by the source
// (other fields are not touched).
func (target *RealmMemberAttrs) Init(source map[string]interface{}) {
	if x := source["operator"]; x != nil {
		target.Operator = true
	}
}

// TagAttrs.  https://ninchat.com/api/v2#tag
type TagAttrs struct {
	Name     *string       `json:"name,omitempty"`
	ParentId *string       `json:"parent_id,omitempty"`
	Theme    *TagThemeAttr `json:"theme,omitempty"`
}

// NewTagAttrs creates an object with the attributes specified by the source.
func NewTagAttrs(source map[string]interface{}) (target *TagAttrs) {
	target = new(TagAttrs)
	target.Init(source)
	return
}

// Init fills in the attributes specified by the source
// (other fields are not touched).
func (target *TagAttrs) Init(source map[string]interface{}) {
	if x := source["name"]; x != nil {
		if y, ok := x.(string); ok {
			target.Name = &y
		}
	}

	if x := source["parent_id"]; x != nil {
		if y, ok := x.(string); ok {
			target.ParentId = &y
		}
	}

	if x := source["theme"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.Theme = NewTagThemeAttr(y)
		}
	}
}

// UserAttrs.  https://ninchat.com/api/v2#user
type UserAttrs struct {
	Admin     bool          `json:"admin,omitempty"`
	Connected bool          `json:"connected,omitempty"`
	Deleted   bool          `json:"deleted,omitempty"`
	Guest     bool          `json:"guest,omitempty"`
	Iconurl   *string       `json:"iconurl,omitempty"`
	Idle      *int          `json:"idle,omitempty"`
	Info      *UserInfoAttr `json:"info,omitempty"`
	Name      *string       `json:"name,omitempty"`
	Realname  *string       `json:"realname,omitempty"`
}

// NewUserAttrs creates an object with the attributes specified by the source.
func NewUserAttrs(source map[string]interface{}) (target *UserAttrs) {
	target = new(UserAttrs)
	target.Init(source)
	return
}

// Init fills in the attributes specified by the source
// (other fields are not touched).
func (target *UserAttrs) Init(source map[string]interface{}) {
	if x := source["admin"]; x != nil {
		target.Admin = true
	}

	if x := source["connected"]; x != nil {
		target.Connected = true
	}

	if x := source["deleted"]; x != nil {
		target.Deleted = true
	}

	if x := source["guest"]; x != nil {
		target.Guest = true
	}

	if x := source["iconurl"]; x != nil {
		if y, ok := x.(string); ok {
			target.Iconurl = &y
		}
	}

	if x := source["idle"]; x != nil {
		if y, ok := x.(int); ok {
			target.Idle = &y
		}
	}

	if x := source["info"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.Info = NewUserInfoAttr(y)
		}
	}

	if x := source["name"]; x != nil {
		if y, ok := x.(string); ok {
			target.Name = &y
		}
	}

	if x := source["realname"]; x != nil {
		if y, ok := x.(string); ok {
			target.Realname = &y
		}
	}
}
