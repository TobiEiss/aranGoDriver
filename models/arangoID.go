package models

type (
	// ArangoID for arangoDB-Identification
	ArangoID struct {
		ID  string `json:"_id"`
		Key string `json:"_key"`
		Rev string `json:"_rev"`
	}
)
