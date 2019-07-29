package aranGoDriver

import (
	"errors"
	"net/http"

	"github.com/TobiEiss/aranGoDriver/aranGoConnection"
	"github.com/TobiEiss/aranGoDriver/models"
)

// AranGoSession represent to Session
type AranGoSession struct {
	arangoCon *aranGoConnection.AranGoConnection
}

const urlAuth = "/_open/auth"
const urlDatabase = "/_api/database"
const urlCollection = "/_api/collection"
const urlDocument = "/_api/document"
const urlCursor = "/_api/cursor"
const urlVersion = "/_api/version"
const urlGraph = "/_api/gharial"
const urlUser = "/_api/user"

const systemDB = "_system"
const migrationColl = "migrations"

// NewAranGoDriverSession creates a new instance of a AranGoDriver-Session.
// Need a host (e.g. "http://localhost:8529/")
func NewAranGoDriverSession(host string) *AranGoSession {
	return &AranGoSession{aranGoConnection.NewAranGoConnection(host)}
}

// Connect to arangoDB
func (session *AranGoSession) Connect(username string, password string) error {
	credentials := models.Credentials{}
	credentials.Username = username
	credentials.Password = password

	var resultMap map[string]string
	err := session.arangoCon.Query(&resultMap, http.MethodPost, urlAuth, credentials)

	if err == nil {
		session.arangoCon.SetJwtKey(resultMap["jwt"])
	}
	return err
}

// Version returns current version
func (session *AranGoSession) Version() (Version, error) {
	var result Version
	err := session.arangoCon.Query(&result, http.MethodGet, urlVersion, nil)
	return result, err
}

// Add a new database user
func (session *AranGoSession) CreateUser(username string, password string) error {
	body := make(map[string]interface{})
	body["user"] = username
	body["passwd"] = password

	var response interface{}

	err := session.arangoCon.Query(&response, http.MethodPost, urlUser, body)

	return err
}

// Delete an existing user
func (session *AranGoSession) DropUser(username string) error {
	var response interface{}

	err := session.arangoCon.Query(&response, http.MethodDelete, urlUser+"/"+username, nil)

	return err
}

// Set the accesslevel for an user on a database
// Possible values for level are: rw, ro and none
func (session *AranGoSession) GrantDB(dbname string, username string, level string) error {
	body := make(map[string]interface{})
	body["grant"] = level

	var response interface{}

	err := session.arangoCon.Query(&response, http.MethodPut, urlUser+"/"+username+"/database/"+dbname, body)

	return err
}

// Set the accesslevel for an user on a collection
// Possible values for level are: rw, ro and none
func (session *AranGoSession) GrantCollection(dbname string, collectionName string, username string, level string) error {
	body := make(map[string]interface{})
	body["grant"] = level

	var response interface{}

	err := session.arangoCon.Query(&response, http.MethodPut, urlUser+"/"+username+"/database/"+dbname+"/"+collectionName, body)

	return err
}

// ListDBs lists all db's
func (session *AranGoSession) ListDBs() ([]string, error) {
	var databaseWrapper struct {
		Databases []string `json:"result,omitempty"`
	}
	err := session.arangoCon.Query(&databaseWrapper, http.MethodGet, urlDatabase, nil)

	return databaseWrapper.Databases, err
}

// CreateDB creates a new db
func (session *AranGoSession) CreateDB(dbname string) error {
	body := make(map[string]string)
	body["name"] = dbname
	var result interface{}
	err := session.arangoCon.Query(&result, http.MethodPost, urlDatabase, body)
	return err
}

// DropDB drop a database
func (session *AranGoSession) DropDB(dbname string) error {
	var response interface{}
	err := session.arangoCon.Query(&response, http.MethodDelete, urlDatabase+"/"+dbname, nil)
	return err
}

// CreateCollection creates a collection
func (session *AranGoSession) CreateCollection(dbname string, collectionName string) error {
	body := make(map[string]string)
	body["name"] = collectionName
	var result interface{}
	err := session.arangoCon.Query(&result, http.MethodPost, "/_db/"+dbname+urlCollection, body)
	return err
}

// CreateEdgeCollection creates a edge to DB
func (session *AranGoSession) CreateEdgeCollection(dbname string, edgeName string) error {
	body := make(map[string]interface{})
	body["name"] = edgeName
	body["type"] = 3
	var result interface{}
	err := session.arangoCon.Query(&result, http.MethodPost, "/_db/"+dbname+urlCollection, body)
	return err
}

func (session *AranGoSession) CreateGraph(dbname string, graphName string, edgeDefinitions []models.EdgeDefinition) error {
	body := make(map[string]interface{})
	body["name"] = graphName
	body["edgeDefinitions"] = edgeDefinitions

	var response interface{}
	err := session.arangoCon.Query(&response, http.MethodPost, "/_db/"+dbname+urlGraph, body)
	return err
}

func (session *AranGoSession) ListGraphs(dbname string) (interface{}, error) {
	var result interface{}
	err := session.arangoCon.Query(&result, http.MethodGet, "/_db/"+dbname+urlGraph, nil)
	return result, err
}

func (session *AranGoSession) DropGraph(dbname string, graphName string) error {
	var response interface{}
	err := session.arangoCon.Query(&response, http.MethodDelete, "/_db/"+dbname+urlGraph+"/"+graphName, nil)
	return err
}

func (session *AranGoSession) CreateEdgeDocument(dbname string, edgeName string, from string, to string) (models.ArangoID, error) {
	body := make(map[string]interface{})
	body["_from"] = from
	body["_to"] = to
	var aranggoID models.ArangoID
	err := session.arangoCon.Query(&aranggoID, http.MethodPost, "/_db/"+dbname+urlDocument+"/"+edgeName, body)
	return aranggoID, err
}

func (session *AranGoSession) ListCollections(dbname string) ([]string, error) {
	var collections []string
	err := session.arangoCon.Query(&collections, http.MethodGet, "/_db/"+dbname+urlCollection, nil)

	return collections, err
}

// DropCollection deletes a collection
func (session *AranGoSession) DropCollection(dbname string, collectionName string) error {
	var response interface{}
	err := session.arangoCon.Query(&response, http.MethodDelete, "/_db/"+dbname+urlCollection+"/"+collectionName, nil)
	return err
}

// TruncateCollection truncate collections
func (session *AranGoSession) TruncateCollection(dbname string, collectionName string) error {
	var response interface{}
	err := session.arangoCon.Query(&response, http.MethodPut, "/_db/"+dbname+urlCollection+"/"+collectionName+"/truncate", nil)
	return err
}

// CreateDocument creates a document in a collection in a database
func (session *AranGoSession) CreateDocument(dbname string, collectionName string, object interface{}) (models.ArangoID, error) {
	var aranggoID models.ArangoID
	err := session.arangoCon.Query(&aranggoID, http.MethodPost, "/_db/"+dbname+urlDocument+"/"+collectionName, object)
	return aranggoID, err
}

// AqlQuery send a query
func (session *AranGoSession) AqlQuery(typ interface{}, dbname string, query string, count bool, batchSize int) error {
	// build request
	requestBody := make(map[string]interface{})
	requestBody["query"] = query
	requestBody["count"] = count
	requestBody["batchSize"] = batchSize

	var result struct {
		Error  bool        `json:"error"`
		Result interface{} `json:"result"`
	}
	result.Result = typ
	err := session.arangoCon.Query(&result, http.MethodPost, "/_db/"+dbname+urlCursor, requestBody)
	if err != nil {
		return err
	}

	if result.Error {
		return errors.New("an error occured")
	}

	return err
}

// GetCollectionByID search collection by id
func (session *AranGoSession) GetCollectionByID(dbname string, id string) (map[string]interface{}, error) {
	var collection map[string]interface{}
	err := session.arangoCon.Query(&collection, http.MethodGet, "/_db/"+dbname+urlCollection+"/"+id, nil)

	return collection, err
}

// UpdateDocument updates an Object
func (session *AranGoSession) UpdateDocument(dbname string, id string, object interface{}) error {
	var result interface{}
	err := session.arangoCon.Query(&result, http.MethodPatch, "/_db/"+dbname+urlDocument+"/"+id, nil)
	return err
}

// Migrate migrates a migration
func (session *AranGoSession) Migrate(migrations ...Migration) error {
	session.CreateCollection(systemDB, migrationColl)

	// helper function
	findMigration := func(name string) (Migration, bool) {
		migrations := []Migration{}
		query := "FOR migration IN " + migrationColl + " FILTER migration.name == '" + name + "' RETURN migration"
		err := session.AqlQuery(&migrations, systemDB, query, true, 1)
		return migrations[0], err == nil
	}

	// iterate all migrations
	for _, mig := range migrations {
		migration, successfully := findMigration(mig.Name)
		if successfully {
			if migration.Status != Finished {
				mig.Handle(session)
				mig.Status = Finished
				session.UpdateDocument(systemDB, mig.ArangoID.ID, mig)
			}
		} else {
			mig.Status = Started
			arangoID, _ := session.CreateDocument(systemDB, migrationColl, mig)
			mig.Handle(session)
			mig.Status = Finished
			session.UpdateDocument(systemDB, arangoID.ID, mig)
		}
	}
	return nil
}
