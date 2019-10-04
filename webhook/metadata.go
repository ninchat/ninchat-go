package webhook

import (
	"encoding/json"
)

// Well-known metadata keys.
const (
	MetadataNinchat = "ninchat"
	MetadataRating  = "rating"
	MetadataSecure  = "secure"
	MetadataTagIDs  = "tag_ids"
	MetadataVars    = "vars"
)

// Well-known secure metadata keys.
const (
	SecureMetadataIdentifier = "identifier"
)

type Metadata map[string]json.RawMessage

func (toplevel Metadata) Ninchat() (ninchat Metadata, err error) {
	raw, found := toplevel[MetadataNinchat]
	if !found {
		return
	}

	err = json.Unmarshal(raw, &ninchat)
	return
}

// Secure metadata unmarshaling may fail e.g. if it's not an object but some
// other type of JSON value.
func (toplevel Metadata) Secure() (secure Metadata, err error) {
	raw, found := toplevel[MetadataSecure]
	if !found {
		return
	}

	err = json.Unmarshal(raw, &secure)
	return
}
