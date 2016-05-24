package ninchatapi

// MemberAttrs action parameter type.
type MemberAttrs interface {
	memberAttrs()
}

func (*ChannelMemberAttrs) memberAttrs() {}
func (*QueueMemberAttrs) memberAttrs()   {}
func (*RealmMemberAttrs) memberAttrs()   {}

// UserInfoAttr represents the user "info" attribute.
type UserInfoAttr struct {
	Company *string `json:"company,omitempty"`
	Url     *string `json:"url,omitempty"`
}

// NewUserInfoAttr.
func NewUserInfoAttr(source map[string]interface{}) (target *UserInfoAttr) {
	target = new(UserInfoAttr)

	if x := source["company"]; x != nil {
		if y, ok := x.(string); ok {
			target.Company = &y
		}
	}

	if x := source["url"]; x != nil {
		if y, ok := x.(string); ok {
			target.Url = &y
		}
	}

	return
}

// RealmOwnerAccountAttr represents the realm "owner_account" attribute.
type RealmOwnerAccountAttr struct {
	Channels     *UserAccountObjects `json:"channels"`
	QueueMembers *UserAccountMembers `json:"queue_members"`
	Queues       *UserAccountObjects `json:"queues"`
}

// NewRealmOwnerAccountAttr.
func NewRealmOwnerAccountAttr(source map[string]interface{}) (target *RealmOwnerAccountAttr) {
	target = new(RealmOwnerAccountAttr)

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

	return
}

// ThemeAttr represents "theme" attributes.
type ThemeAttr struct {
	Color *string `json:"color,omitempty"`
}

func (target *ThemeAttr) init(source map[string]interface{}) {
	if x := source["color"]; x != nil {
		if y, ok := x.(string); ok {
			target.Color = &y
		}
	}
}

// RealmThemeAttr represents the realm "theme" attribute.
type RealmThemeAttr struct {
	ThemeAttr
}

// NewRealmThemeAttr.
func NewRealmThemeAttr(source map[string]interface{}) (target *RealmThemeAttr) {
	target = new(RealmThemeAttr)
	target.init(source)
	return
}

// TagThemeAttr represents the tag "theme" attribute.
type TagThemeAttr struct {
	ThemeAttr
}

// NewTagThemeAttr.
func NewTagThemeAttr(source map[string]interface{}) (target *TagThemeAttr) {
	target = new(TagThemeAttr)
	target.init(source)
	return
}
