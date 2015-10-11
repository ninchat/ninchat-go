package ninchatapi

// UserInfoAttr represents the user "info" attribute.
type UserInfoAttr struct {
	Company *string `json:"company,omitempty"`
	Url     *string `json:"url,omitempty"`
}

// MergeFrom fills in the parameters specified by the source.
func (target *UserInfoAttr) MergeFrom(source map[string]interface{}) {
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
}

// RealmOwnerAccountAttr represents the realm "owner_account" attribute.
type RealmOwnerAccountAttr struct {
	Channels     *UserAccountObjects `json:"channels"`
	QueueMembers *UserAccountMembers `json:"queue_members"`
	Queues       *UserAccountObjects `json:"queues"`
}

// MergeFrom fills in the parameters specified by the source.
func (target *RealmOwnerAccountAttr) MergeFrom(source map[string]interface{}) {
	if x := source["channels"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.Channels = new(UserAccountObjects)
			target.Channels.MergeFrom(y)
		}
	}

	if x := source["queue_members"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.QueueMembers = new(UserAccountMembers)
			target.QueueMembers.MergeFrom(y)
		}
	}

	if x := source["queues"]; x != nil {
		if y, ok := x.(map[string]interface{}); ok {
			target.Queues = new(UserAccountObjects)
			target.Queues.MergeFrom(y)
		}
	}
}

// RealmThemeAttr represents the realm "theme" attribute.
type RealmThemeAttr struct {
	Color *string `json:"color,omitempty"`
}

// MergeFrom fills in the parameters specified by the source.
func (target *RealmThemeAttr) MergeFrom(source map[string]interface{}) {
	if x := source["color"]; x != nil {
		if y, ok := x.(string); ok {
			target.Color = &y
		}
	}
}
