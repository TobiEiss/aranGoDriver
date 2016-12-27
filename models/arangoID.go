package models

type (
	// ArangoID for arangoDB-Identification
	ArangoID struct {
		ID  string `json:"_id,omitempty"`
		Key string `json:"_key,omitempty"`
		Rev string `json:"_rev,omitempty"`
	}
)
