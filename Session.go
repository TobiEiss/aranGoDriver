package aranGoDriver

import (
	"github.com/TobiEiss/aranGoDriver/models"
)

type Session interface {
	Connect(username string, password string) error
	Version() (Version, error)

	// users
	CreateUser(username string, password string) error
	DropUser(username string) error
	GrantDB(dbname string, username string, level string) error
	GrantCollection(dbname string, collectionName string, username string, level string) error

	// databases
	ListDBs() ([]string, error)
	CreateDB(dbname string) error
	DropDB(dbname string) error

	ListCollections(dbname string) ([]string, error)
	CreateCollection(dbname string, collectionName string) error
	DropCollection(dbname string, collectionName string) error
	TruncateCollection(dbname string, collectionName string) error

	CreateEdgeCollection(dbname string, edgeName string) error
	CreateEdgeDocument(dbname string, edgeName string, from string, to string) (models.ArangoID, error)

	CreateGraph(dbname string, graphName string, edgeDefinitions []models.EdgeDefinition) error
	ListGraphs(dbname string) (interface{}, error)
	DropGraph(dbname string, graphName string) error

	// GetCollectionByID search collection by id
	// returns:
	// -> result as map
	// -> error if applicable
	GetCollectionByID(dbname string, id string) (map[string]interface{}, error)
	CreateDocument(dbname string, collectionName string, object interface{}) (models.ArangoID, error)
	UpdateDocument(dbname string, id string, object interface{}) error

	// AqlQuery returns: result as array-map, error
	AqlQuery(typ interface{}, dbname string, query string, count bool, batchSize int) error

	// Query with auth
	Query(typ interface{}, methode string, route string, body interface{}) error

	Migrate(migration ...Migration) error
}

// Version represent the version and license of the ArangoDB
type Version struct {
	Server  string `json:"server"`
	License string `json:"license"`
}
