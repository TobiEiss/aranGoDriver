package models

type (
	// Credentials for arangoDB
	Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
)
