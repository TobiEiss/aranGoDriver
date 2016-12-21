package aranGoDriver

import (
	"github.com/TobiEiss/aranGoDriver/models"
)

type Session interface {
	Connect(username string, password string) error

	// databases
	ListDBs() ([]string, error)
	CreateDB(dbname string) error
	DropDB(dbname string) error

	CreateCollection(dbname string, collectionName string) error
	DropCollection(dbname string, collectionName string) error
	TruncateCollection(dbname string, collectionName string) error

	// GetCollectionByID search collection by id
	// returns:
	// -> result as jsonString
	// -> result as map
	// -> error if applicable
	GetCollectionByID(dbname string, id string) (string, map[string]interface{}, error)

	CreateDocument(dbname string, collectionName string, object map[string]interface{}) (models.ArangoID, error)
	// AqlQuery returns: result as array-map, result as json, error
	AqlQuery(dbname string, query string, count bool, batchSize int) ([]map[string]interface{}, string, error)
}
