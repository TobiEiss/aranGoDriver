package framed

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
	"github.com/TobiEiss/aranGoDriver"
)

// Database represent a arangodb-database
type Database struct {
	Name    string
	Session *aranGoDriver.Session
}

// ListCollections returns an array of all collections in catabase
func (database *Database) ListCollections() ([]Collection, error) {
	collections := make([]Collection, 0)

	// request collections
	response, _, err := (*database.Session).ListCollections(database.Name)
	jsonParsed, err := gabs.ParseJSON([]byte(response))
	if err != nil {
		return collections, err
	}

	// parse collections
	children, _ := jsonParsed.S("result").Children()
	for _, child := range children {
		var collection Collection
		json.Unmarshal(child.Bytes(), &collection)
		collection.Database = database
		collections = append(collections, collection)
	}
	return collections, nil
}
