package ninchatapi

// THIS FILE IS AUTO-GENERATED BY generate.py - DO NOT EDIT BY HAND!

// ChannelMember event parameter type.
type ChannelMember struct {
	MemberAttrs *ChannelMemberAttrs `json:"member_attrs"`
	UserAttrs   *UserAttrs          `json:"user_attrs"`
}

// NewChannelMember creates an object with the parameters specified by the source.
func NewChannelMember(source map[string]interface{}) (target *ChannelMember) {
	target = new(ChannelMember)
	target.Init(source)
	return
}

// Init fills in the parameters specified by the source
// (other fields are not touched).
func (target *ChannelMember) Init(source map[string]interface{}) {
	if x := source["member_attrs"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.MemberAttrs = NewChannelMemberAttrs(y)
		}
	}

	if x := source["user_attrs"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.UserAttrs = NewUserAttrs(y)
		}
	}
}

// MakeChannelMembers duplicates the map while unwrapping the values.
func MakeChannelMembers(source map[string]interface{}) (target map[string]*ChannelMember) {
	target = make(map[string]*ChannelMember)

	for key, x := range source {
		if y, ok := x.(map[string]interface{}); ok {
			target[key] = NewChannelMember(y)
		}
	}

	return
}

// ChannelResult event parameter type.
type ChannelResult struct {
	ChannelAttrs *ChannelAttrs `json:"channel_attrs"`
	RealmId      *string       `json:"realm_id,omitempty"`
	Weight       float64       `json:"weight"`
}

// NewChannelResult creates an object with the parameters specified by the source.
func NewChannelResult(source map[string]interface{}) (target *ChannelResult) {
	target = new(ChannelResult)
	target.Init(source)
	return
}

// Init fills in the parameters specified by the source
// (other fields are not touched).
func (target *ChannelResult) Init(source map[string]interface{}) {
	if x := source["channel_attrs"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.ChannelAttrs = NewChannelAttrs(y)
		}
	}

	if x := source["realm_id"]; x != nil {
		if y, ok := x.(string); ok {
			target.RealmId = &y
		}
	}

	if x := source["weight"]; x != nil {
		if y, ok := x.(float64); ok {
			target.Weight = y
		}
	}
}

// MakeChannels duplicates the map while unwrapping the values.
func MakeChannels(source map[string]interface{}) (target map[string]*ChannelResult) {
	target = make(map[string]*ChannelResult)

	for key, x := range source {
		if y, ok := x.(map[string]interface{}); ok {
			target[key] = NewChannelResult(y)
		}
	}

	return
}

// MakeDialogueMembers duplicates the map while unwrapping the values.
func MakeDialogueMembers(source map[string]interface{}) (target map[string]*DialogueMemberAttrs) {
	target = make(map[string]*DialogueMemberAttrs)

	for key, x := range source {
		if y, ok := x.(map[string]interface{}); ok {
			target[key] = NewDialogueMemberAttrs(y)
		}
	}

	return
}

// MasterKey event parameter type.
type MasterKey struct {
}

// NewMasterKey creates an object with the parameters specified by the source.
func NewMasterKey(source map[string]interface{}) (target *MasterKey) {
	target = new(MasterKey)
	target.Init(source)
	return
}

// Init fills in the parameters specified by the source
// (other fields are not touched).
func (target *MasterKey) Init(source map[string]interface{}) {
}

// MakeMasterKeys duplicates the map while unwrapping the values.
func MakeMasterKeys(source map[string]interface{}) (target map[string]*MasterKey) {
	target = make(map[string]*MasterKey)

	for key, x := range source {
		if y, ok := x.(map[string]interface{}); ok {
			target[key] = NewMasterKey(y)
		}
	}

	return
}

// QueueMember event parameter type.
type QueueMember struct {
	MemberAttrs *QueueMemberAttrs `json:"member_attrs"`
	UserAttrs   *UserAttrs        `json:"user_attrs"`
}

// NewQueueMember creates an object with the parameters specified by the source.
func NewQueueMember(source map[string]interface{}) (target *QueueMember) {
	target = new(QueueMember)
	target.Init(source)
	return
}

// Init fills in the parameters specified by the source
// (other fields are not touched).
func (target *QueueMember) Init(source map[string]interface{}) {
	if x := source["member_attrs"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.MemberAttrs = NewQueueMemberAttrs(y)
		}
	}

	if x := source["user_attrs"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.UserAttrs = NewUserAttrs(y)
		}
	}
}

// MakeQueueMembers duplicates the map while unwrapping the values.
func MakeQueueMembers(source map[string]interface{}) (target map[string]*QueueMember) {
	target = make(map[string]*QueueMember)

	for key, x := range source {
		if y, ok := x.(map[string]interface{}); ok {
			target[key] = NewQueueMember(y)
		}
	}

	return
}

// RealmMember event parameter type.
type RealmMember struct {
	MemberAttrs *RealmMemberAttrs `json:"member_attrs"`
	UserAttrs   *UserAttrs        `json:"user_attrs"`
}

// NewRealmMember creates an object with the parameters specified by the source.
func NewRealmMember(source map[string]interface{}) (target *RealmMember) {
	target = new(RealmMember)
	target.Init(source)
	return
}

// Init fills in the parameters specified by the source
// (other fields are not touched).
func (target *RealmMember) Init(source map[string]interface{}) {
	if x := source["member_attrs"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.MemberAttrs = NewRealmMemberAttrs(y)
		}
	}

	if x := source["user_attrs"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.UserAttrs = NewUserAttrs(y)
		}
	}
}

// MakeRealmMembers duplicates the map while unwrapping the values.
func MakeRealmMembers(source map[string]interface{}) (target map[string]*RealmMember) {
	target = make(map[string]*RealmMember)

	for key, x := range source {
		if y, ok := x.(map[string]interface{}); ok {
			target[key] = NewRealmMember(y)
		}
	}

	return
}

// RealmQueue event parameter type.
type RealmQueue struct {
	QueueAttrs    *QueueAttrs `json:"queue_attrs"`
	QueuePosition *int        `json:"queue_position,omitempty"`
}

// NewRealmQueue creates an object with the parameters specified by the source.
func NewRealmQueue(source map[string]interface{}) (target *RealmQueue) {
	target = new(RealmQueue)
	target.Init(source)
	return
}

// Init fills in the parameters specified by the source
// (other fields are not touched).
func (target *RealmQueue) Init(source map[string]interface{}) {
	if x := source["queue_attrs"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.QueueAttrs = NewQueueAttrs(y)
		}
	}

	if x := source["queue_position"]; x != nil {
		if y, ok := x.(int); ok {
			target.QueuePosition = &y
		}
	}
}

// MakeRealmQueues duplicates the map while unwrapping the values.
func MakeRealmQueues(source map[string]interface{}) (target map[string]*RealmQueue) {
	target = make(map[string]*RealmQueue)

	for key, x := range source {
		if y, ok := x.(map[string]interface{}); ok {
			target[key] = NewRealmQueue(y)
		}
	}

	return
}

// MakeStrings duplicates the map while unwrapping the values.
func MakeStrings(source map[string]interface{}) (target map[string]string) {
	target = make(map[string]string)

	for key, x := range source {
		if y, ok := x.(string); ok {
			target[key] = y
		}
	}

	return
}

// UserAccount event parameter type.
type UserAccount struct {
	Channels      *UserAccountObjects        `json:"channels,omitempty"`
	QueueMembers  *UserAccountMembers        `json:"queue_members,omitempty"`
	Queues        *UserAccountObjects        `json:"queues,omitempty"`
	Realms        *UserAccountObjects        `json:"realms,omitempty"`
	Subscriptions []*UserAccountSubscription `json:"subscriptions,omitempty"`
}

// NewUserAccount creates an object with the parameters specified by the source.
func NewUserAccount(source map[string]interface{}) (target *UserAccount) {
	target = new(UserAccount)
	target.Init(source)
	return
}

// Init fills in the parameters specified by the source
// (other fields are not touched).
func (target *UserAccount) Init(source map[string]interface{}) {
	if x := source["channels"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.Channels = NewUserAccountObjects(y)
		}
	}

	if x := source["queue_members"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.QueueMembers = NewUserAccountMembers(y)
		}
	}

	if x := source["queues"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.Queues = NewUserAccountObjects(y)
		}
	}

	if x := source["realms"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.Realms = NewUserAccountObjects(y)
		}
	}

	if x := source["subscriptions"]; x != nil {
		if y, ok := x.([]interface{}); ok {
			target.Subscriptions = AppendUserAccountSubscriptions(nil, y)
		}
	}
}

// UserAccountMembers event parameter type.
type UserAccountMembers struct {
	Quota int `json:"quota"`
}

// NewUserAccountMembers creates an object with the parameters specified by the source.
func NewUserAccountMembers(source map[string]interface{}) (target *UserAccountMembers) {
	target = new(UserAccountMembers)
	target.Init(source)
	return
}

// Init fills in the parameters specified by the source
// (other fields are not touched).
func (target *UserAccountMembers) Init(source map[string]interface{}) {
	if x := source["quota"]; x != nil {
		if y, ok := x.(int); ok {
			target.Quota = y
		}
	}
}

// UserAccountObjects event parameter type.
type UserAccountObjects struct {
	Available int `json:"available"`
	Quota     int `json:"quota"`
}

// NewUserAccountObjects creates an object with the parameters specified by the source.
func NewUserAccountObjects(source map[string]interface{}) (target *UserAccountObjects) {
	target = new(UserAccountObjects)
	target.Init(source)
	return
}

// Init fills in the parameters specified by the source
// (other fields are not touched).
func (target *UserAccountObjects) Init(source map[string]interface{}) {
	if x := source["available"]; x != nil {
		if y, ok := x.(int); ok {
			target.Available = y
		}
	}

	if x := source["quota"]; x != nil {
		if y, ok := x.(int); ok {
			target.Quota = y
		}
	}
}

// UserAccountSubscription event parameter type.
type UserAccountSubscription struct {
	Active       bool                `json:"active"`
	Channels     *UserAccountObjects `json:"channels,omitempty"`
	Expiration   *int                `json:"expiration,omitempty"`
	Plan         string              `json:"plan"`
	QueueMembers *UserAccountMembers `json:"queue_members,omitempty"`
	Queues       *UserAccountObjects `json:"queues,omitempty"`
	Realms       *UserAccountObjects `json:"realms,omitempty"`
	Renewal      *int                `json:"renewal,omitempty"`
}

// NewUserAccountSubscription creates an object with the parameters specified by the source.
func NewUserAccountSubscription(source map[string]interface{}) (target *UserAccountSubscription) {
	target = new(UserAccountSubscription)
	target.Init(source)
	return
}

// Init fills in the parameters specified by the source
// (other fields are not touched).
func (target *UserAccountSubscription) Init(source map[string]interface{}) {
	if x := source["active"]; x != nil {
		target.Active = true
	}

	if x := source["channels"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.Channels = NewUserAccountObjects(y)
		}
	}

	if x := source["expiration"]; x != nil {
		if y, ok := x.(int); ok {
			target.Expiration = &y
		}
	}

	if x := source["plan"]; x != nil {
		if y, ok := x.(string); ok {
			target.Plan = y
		}
	}

	if x := source["queue_members"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.QueueMembers = NewUserAccountMembers(y)
		}
	}

	if x := source["queues"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.Queues = NewUserAccountObjects(y)
		}
	}

	if x := source["realms"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.Realms = NewUserAccountObjects(y)
		}
	}

	if x := source["renewal"]; x != nil {
		if y, ok := x.(int); ok {
			target.Renewal = &y
		}
	}
}

// AppendUserAccountSubscriptions duplicates the source slice while unwrapping the elements.
func AppendUserAccountSubscriptions(target []*UserAccountSubscription, source []interface{}) []*UserAccountSubscription {
	if source != nil {
		if target == nil || cap(target) < len(target)+len(source) {
			t := make([]*UserAccountSubscription, len(target), len(target)+len(source))
			copy(t, target)
			target = t
		}

		for _, x := range source {
			var z *UserAccountSubscription
			if y, ok := x.(map[string]interface{}); ok {
				z = NewUserAccountSubscription(y)
			}
			target = append(target, z)
		}
	}

	return target
}

// UserChannel event parameter type.
type UserChannel struct {
	ChannelAttrs  *ChannelAttrs `json:"channel_attrs"`
	ChannelStatus *string       `json:"channel_status,omitempty"`
	RealmId       *string       `json:"realm_id,omitempty"`
}

// NewUserChannel creates an object with the parameters specified by the source.
func NewUserChannel(source map[string]interface{}) (target *UserChannel) {
	target = new(UserChannel)
	target.Init(source)
	return
}

// Init fills in the parameters specified by the source
// (other fields are not touched).
func (target *UserChannel) Init(source map[string]interface{}) {
	if x := source["channel_attrs"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.ChannelAttrs = NewChannelAttrs(y)
		}
	}

	if x := source["channel_status"]; x != nil {
		if y, ok := x.(string); ok {
			target.ChannelStatus = &y
		}
	}

	if x := source["realm_id"]; x != nil {
		if y, ok := x.(string); ok {
			target.RealmId = &y
		}
	}
}

// MakeUserChannels duplicates the map while unwrapping the values.
func MakeUserChannels(source map[string]interface{}) (target map[string]*UserChannel) {
	target = make(map[string]*UserChannel)

	for key, x := range source {
		if y, ok := x.(map[string]interface{}); ok {
			target[key] = NewUserChannel(y)
		}
	}

	return
}

// UserDialogue event parameter type.
type UserDialogue struct {
	AudienceMetadata map[string]string               `json:"audience_metadata,omitempty"`
	DialogueMembers  map[string]*DialogueMemberAttrs `json:"dialogue_members,omitempty"`
	DialogueStatus   *string                         `json:"dialogue_status,omitempty"`
}

// NewUserDialogue creates an object with the parameters specified by the source.
func NewUserDialogue(source map[string]interface{}) (target *UserDialogue) {
	target = new(UserDialogue)
	target.Init(source)
	return
}

// Init fills in the parameters specified by the source
// (other fields are not touched).
func (target *UserDialogue) Init(source map[string]interface{}) {
	if x := source["audience_metadata"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.AudienceMetadata = MakeStrings(y)
		}
	}

	if x := source["dialogue_members"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.DialogueMembers = MakeDialogueMembers(y)
		}
	}

	if x := source["dialogue_status"]; x != nil {
		if y, ok := x.(string); ok {
			target.DialogueStatus = &y
		}
	}
}

// MakeUserDialogues duplicates the map while unwrapping the values.
func MakeUserDialogues(source map[string]interface{}) (target map[string]*UserDialogue) {
	target = make(map[string]*UserDialogue)

	for key, x := range source {
		if y, ok := x.(map[string]interface{}); ok {
			target[key] = NewUserDialogue(y)
		}
	}

	return
}

// MakeUserIdentities duplicates the map while unwrapping the values.
func MakeUserIdentities(source map[string]interface{}) (target map[string]*IdentityAttrs) {
	target = make(map[string]*IdentityAttrs)

	for key, x := range source {
		if y, ok := x.(map[string]interface{}); ok {
			target[key] = NewIdentityAttrs(y)
		}
	}

	return
}

// UserQueue event parameter type.
type UserQueue struct {
	QueueAttrs *QueueAttrs `json:"queue_attrs"`
	RealmId    string      `json:"realm_id"`
}

// NewUserQueue creates an object with the parameters specified by the source.
func NewUserQueue(source map[string]interface{}) (target *UserQueue) {
	target = new(UserQueue)
	target.Init(source)
	return
}

// Init fills in the parameters specified by the source
// (other fields are not touched).
func (target *UserQueue) Init(source map[string]interface{}) {
	if x := source["queue_attrs"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.QueueAttrs = NewQueueAttrs(y)
		}
	}

	if x := source["realm_id"]; x != nil {
		if y, ok := x.(string); ok {
			target.RealmId = y
		}
	}
}

// MakeUserQueues duplicates the map while unwrapping the values.
func MakeUserQueues(source map[string]interface{}) (target map[string]*UserQueue) {
	target = make(map[string]*UserQueue)

	for key, x := range source {
		if y, ok := x.(map[string]interface{}); ok {
			target[key] = NewUserQueue(y)
		}
	}

	return
}

// MakeUserRealms duplicates the map while unwrapping the values.
func MakeUserRealms(source map[string]interface{}) (target map[string]*RealmAttrs) {
	target = make(map[string]*RealmAttrs)

	for key, x := range source {
		if y, ok := x.(map[string]interface{}); ok {
			target[key] = NewRealmAttrs(y)
		}
	}

	return
}

// MakeUserRealmsMember duplicates the map while unwrapping the values.
func MakeUserRealmsMember(source map[string]interface{}) (target map[string]*RealmMemberAttrs) {
	target = make(map[string]*RealmMemberAttrs)

	for key, x := range source {
		if y, ok := x.(map[string]interface{}); ok {
			target[key] = NewRealmMemberAttrs(y)
		}
	}

	return
}

// UserResult event parameter type.
type UserResult struct {
	UserAttrs *UserAttrs `json:"user_attrs"`
	Weight    float64    `json:"weight"`
}

// NewUserResult creates an object with the parameters specified by the source.
func NewUserResult(source map[string]interface{}) (target *UserResult) {
	target = new(UserResult)
	target.Init(source)
	return
}

// Init fills in the parameters specified by the source
// (other fields are not touched).
func (target *UserResult) Init(source map[string]interface{}) {
	if x := source["user_attrs"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.UserAttrs = NewUserAttrs(y)
		}
	}

	if x := source["weight"]; x != nil {
		if y, ok := x.(float64); ok {
			target.Weight = y
		}
	}
}

// MakeUsers duplicates the map while unwrapping the values.
func MakeUsers(source map[string]interface{}) (target map[string]*UserResult) {
	target = make(map[string]*UserResult)

	for key, x := range source {
		if y, ok := x.(map[string]interface{}); ok {
			target[key] = NewUserResult(y)
		}
	}

	return
}