package aranGoDriver

import (
	"encoding/json"

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

	_, resp, err := session.arangoCon.Post(urlAuth, credentials)
	session.arangoCon.SetJwtKey(resp["jwt"].(string))
	return err
}

// ListDBs lists all db's
func (session *AranGoSession) ListDBs() ([]string, error) {
	_, resp, err := session.arangoCon.Get(urlDatabase)
	result := resp["result"].([]interface{})

	dblist := make([]string, len(result))
	for _, value := range result {
		if str, ok := value.(string); ok {
			dblist = append(dblist, str)
		}
	}

	return dblist, err
}

// CreateDB creates a new db
func (session *AranGoSession) CreateDB(dbname string) error {
	body := make(map[string]string)
	body["name"] = dbname

	_, _, err := session.arangoCon.Post(urlDatabase, body)
	return err
}

// DropDB drop a database
func (session *AranGoSession) DropDB(dbname string) error {
	_, _, err := session.arangoCon.Delete(urlDatabase + "/" + dbname)
	return err
}

// CreateCollection creates a collection
func (session *AranGoSession) CreateCollection(dbname string, collectionName string) error {
	body := make(map[string]string)
	body["name"] = collectionName
	_, _, err := session.arangoCon.Post("/_db/"+dbname+urlCollection, body)
	return err
}

// DropCollection deletes a collection
func (session *AranGoSession) DropCollection(dbname string, collectionName string) error {
	_, _, err := session.arangoCon.Delete("/_db/" + dbname + urlCollection + "/" + collectionName)
	return err
}

// TruncateCollection truncate collections
func (session *AranGoSession) TruncateCollection(dbname string, collectionName string) error {
	_, _, err := session.arangoCon.Put("/_db/"+dbname+urlCollection+"/"+collectionName+"/truncate", "")
	return err
}

// CreateDocument creates a document in a dollection in a database
func (session *AranGoSession) CreateDocument(dbname string, collectionName string, object map[string]interface{}) (models.ArangoID, error) {
	bodyString, _, err := session.arangoCon.Post("/_db/"+dbname+urlDocument+"/"+collectionName, object)
	aranggoID := models.ArangoID{}
	err = json.Unmarshal([]byte(bodyString), &aranggoID)
	return aranggoID, err
}

// AqlQuery send a query
func (session *AranGoSession) AqlQuery(dbname string, query string, count bool, batchSize int) (map[string]interface{}, error) {
	requestBody := make(map[string]interface{})
	requestBody["query"] = query
	requestBody["count"] = count
	requestBody["batchSize"] = batchSize
	_, result, err := session.arangoCon.Post("/_db/"+dbname+urlCursor, requestBody)
	return result, err
}
