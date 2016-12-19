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

	CreateDocument(dbname string, collectionName string, object map[string]interface{}) (models.ArangoID, error)
	// AqlQuery returns: result as array-map, result as json, error
	AqlQuery(dbname string, query string, count bool, batchSize int) ([]map[string]interface{}, string, error)
}
