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

// CreateCollection creates a new collection.
// After that search for the right collection via ListCollections().
func (database *Database) CreateCollection(collectionname string) Collection {
	var collection Collection

	// create new collection
	err := (*database.Session).CreateCollection(database.Name, collectionname)
	if err != nil {
		return collection
	}

	// get all collections
	allCollections, err := database.ListCollections()
	if err != nil {
		return collection
	}
	// try to find the new collection
	for _, c := range allCollections {
		if c.Name == collectionname {
			collection = c
		}
	}

	return collection
}
