package ninchatapi

// UserInfoAttr represents the user "info" attribute.
type UserInfoAttr struct {
	Company *string `json:"company,omitempty"`
	Url     *string `json:"url,omitempty"`
}

func (target *UserInfoAttr) init(source map[string]interface{}) {
	// TODO
}

// RealmOwnerAccountAttr represents the realm "owner_account" attribute.
type RealmOwnerAccountAttr struct {
	Channels     *UserAccountObjects `json:"channels"`
	QueueMembers *UserAccountMembers `json:"queue_members"`
	Queues       *UserAccountObjects `json:"queues"`
}

func (target *RealmOwnerAccountAttr) init(source map[string]interface{}) {
	// TODO
}

// RealmThemeAttr represents the realm "theme" attribute.
type RealmThemeAttr struct {
	Color *string `json:"color,omitempty"`
}

func (target *RealmThemeAttr) init(source map[string]interface{}) {
	// TODO
}
