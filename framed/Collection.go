package framed

// Collection represent an arango-collection
type Collection struct {
	Name     string `json:"name,omitempty"`
	IsSystem bool   `json:"isSystem,omitempty"`
	Status   int    `json:"status,omitempty"`
	Type     int    `json:"type,omitempty"`
	Database *Database
}
