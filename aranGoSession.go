package aranGoDriver

import (
	"encoding/json"
	"net/http"
	"reflect"

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

	_, resp, err := session.arangoCon.Post(urlAuth, credentials)
	if err == nil {
		session.arangoCon.SetJwtKey(resp["jwt"].(string))
	}
	return err
}

// ListDBs lists all db's
func (session *AranGoSession) ListDBs() ([]string, error) {
	var databaseWrapper struct {
		Databases []string `json:"result,omitempty"`
	}
	err := session.arangoCon.Query(&databaseWrapper, http.MethodGet, urlDatabase)

	return databaseWrapper.Databases, err
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

// CreateEdgeCollection creates a edge to DB
func (session *AranGoSession) CreateEdgeCollection(dbname string, edgeName string) error {
	body := make(map[string]interface{})
	body["name"] = edgeName
	body["type"] = 3

	_, _, err := session.arangoCon.Post("/_db/"+dbname+urlCollection, body)
	return err
}

func (session *AranGoSession) CreateEdgeDocument(dbname string, edgeName string, from string, to string) (models.ArangoID, error) {
	body := make(map[string]interface{})
	body["_from"] = from
	body["_to"] = to

	bodyString, _, err := session.arangoCon.Post("/_db/"+dbname+urlDocument+"/"+edgeName, body)
	aranggoID := models.ArangoID{}
	err = json.Unmarshal([]byte(bodyString), &aranggoID)
	return aranggoID, err
}

func (session *AranGoSession) ListCollections(dbname string) (map[string]interface{}, error) {
	var collections map[string]interface{}
	err := session.arangoCon.Query(&collections, http.MethodGet, "/_db/"+dbname+urlCollection)

	return collections, err
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

// CreateDocument creates a document in a collection in a database
func (session *AranGoSession) CreateDocument(dbname string, collectionName string, object interface{}) (models.ArangoID, error) {
	bodyString, _, err := session.arangoCon.Post("/_db/"+dbname+urlDocument+"/"+collectionName, object)
	aranggoID := models.ArangoID{}
	err = json.Unmarshal([]byte(bodyString), &aranggoID)
	return aranggoID, err
}

// CreateJSONDocument creates a document in a collection in a database
func (session *AranGoSession) CreateJSONDocument(dbname string, collectionName string, jsonObj string) (models.ArangoID, error) {
	bodyString, _, err := session.arangoCon.PostJSON("/_db/"+dbname+urlDocument+"/"+collectionName, []byte(jsonObj))
	aranggoID := models.ArangoID{}
	err = json.Unmarshal([]byte(bodyString), &aranggoID)
	return aranggoID, err
}

// AqlQuery send a query
func (session *AranGoSession) AqlQuery(dbname string, query string, count bool, batchSize int) ([]map[string]interface{}, string, error) {
	requestBody := make(map[string]interface{})
	requestBody["query"] = query
	requestBody["count"] = count
	requestBody["batchSize"] = batchSize
	_, response, err := session.arangoCon.Post("/_db/"+dbname+urlCursor, requestBody)
	if err != nil {
		return nil, "", err
	}

	// map response to array of map
	resultInterface := response["result"]
	resultSlice := reflect.ValueOf(resultInterface)

	if errorBool := response["error"].(bool); !errorBool && resultSlice.Len() > 0 {
		result := make([]map[string]interface{}, resultSlice.Len())
		for i := 0; i < resultSlice.Len(); i++ {
			result[i] = resultSlice.Index(i).Interface().(map[string]interface{})
		}

		// only result as json
		resultByte, err := json.Marshal(resultInterface)

		return result, string(resultByte), err
	}
	return nil, "", err
}

// GetCollectionByID search collection by id
func (session *AranGoSession) GetCollectionByID(dbname string, id string) (map[string]interface{}, error) {
	var collection map[string]interface{}
	err := session.arangoCon.Query(&collection, http.MethodGet, "/_db/"+dbname+urlDocument+"/"+id)

	return collection, err
}

// UpdateDocument updates an Object
func (session *AranGoSession) UpdateDocument(dbname string, id string, object interface{}) error {
	_, _, err := session.arangoCon.Patch("/_db/"+dbname+urlDocument+"/"+id, object)
	return err
}

// UpdateJSONDocument update a json
func (session *AranGoSession) UpdateJSONDocument(dbname string, id string, jsonObj string) error {
	_, _, err := session.arangoCon.PatchJSON("/_db/"+dbname+urlDocument+"/"+id, []byte(jsonObj))
	return err
}

// Migrate migrates a migration
func (session *AranGoSession) Migrate(migrations ...Migration) error {
	session.CreateCollection(systemDB, migrationColl)

	// helper function
	findMigration := func(name string) (Migration, bool) {
		query := "FOR migration IN " + migrationColl + " FILTER migration.name == '" + name + "' RETURN migration"
		_, jsonMig, err := session.AqlQuery(systemDB, query, true, 1)
		migrations := []Migration{}
		err = json.Unmarshal([]byte(jsonMig), &migrations)
		if jsonMig == "" || err != nil {
			return Migration{}, false
		}
		return migrations[0], jsonMig != "" && err == nil
	}

	migrationToJSON := func(migration Migration) string {
		b, _ := json.Marshal(migration)
		return string(b)
	}

	// iterate all migrations
	for _, mig := range migrations {
		migration, successfully := findMigration(mig.Name)
		if successfully {
			if migration.Status != Finished {
				mig.Handle(session)
				mig.Status = Finished
				session.UpdateJSONDocument(systemDB, mig.ArangoID.ID, migrationToJSON(mig))
			}
		} else {
			mig.Status = Started
			arangoID, _ := session.CreateJSONDocument(systemDB, migrationColl, migrationToJSON(mig))
			mig.Handle(session)
			mig.Status = Finished
			session.UpdateJSONDocument(systemDB, arangoID.ID, migrationToJSON(mig))
		}
	}
	return nil
}
