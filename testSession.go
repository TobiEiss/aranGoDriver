package aranGoDriver

import (
	"fmt"
	"strconv"
	"time"

	"github.com/TobiEiss/aranGoDriver/models"

	"errors"
	"math/rand"

	"encoding/json"

	"github.com/TobiEiss/aranGoDriver/sliceTricks"
	"github.com/fatih/structs"
)

type TestSession struct {
	database map[string]map[string][]map[string]interface{}
	aqlFakes map[string]AqlFake
}

type AqlFake struct {
	JsonResult string
	MapResult  []map[string]interface{}
}

func NewTestSession() *TestSession {
	// database - collection - list of document (key, value)
	testSession := &TestSession{make(map[string]map[string][]map[string]interface{}), make(map[string]AqlFake)}
	testSession.database[systemDB] = make(map[string][]map[string]interface{})
	return testSession
}

func findByParam(session *TestSession, dbname string, keyName string, valueV string) *map[string]interface{} {
	for _, collection := range session.database[dbname] {
		for _, entry := range collection {
			for key, value := range entry {
				if key == keyName && value == valueV {
					return &entry
				}
			}
		}
	}
	return nil
}

func objectToMap(object interface{}) map[string]interface{} {
	var entryAsMap map[string]interface{}
	switch item := object.(type) {
	case map[string]interface{}:
		entryAsMap = item
	default:
		entryAsMap = structs.Map(object)
	}
	return entryAsMap
}

// Connect test
func (session TestSession) Connect(username string, password string) error {
	fmt.Println("Connect to DB")
	return nil
}

func (session *TestSession) ListDBs() ([]string, error) {
	databases := []string{}

	for key := range session.database {
		databases = append(databases, key)
	}

	return databases, nil
}

// CreateDB test create a db
func (session *TestSession) CreateDB(dbname string) error {
	_, ok := session.database[dbname]
	if ok {
		return errors.New("DB already exists")
	}
	session.database[dbname] = make(map[string][]map[string]interface{})
	return nil
}

func (session *TestSession) DropDB(dbname string) error {
	delete(session.database, dbname)
	return nil
}

func (session *TestSession) CreateCollection(dbname string, collectionName string) error {
	_, ok := session.database[dbname]
	if !ok {
		return errors.New("DB doesnt")
	}

	session.database[dbname][collectionName] = make([]map[string]interface{}, 10)
	return nil
}

func (session *TestSession) CreateEdgeCollection(dbname string, edgeName string) error {
	return errors.New("not implemented..")
}

func (session *TestSession) CreateEdgeDocument(dbname string, edgeName string, from string, to string) (models.ArangoID, error) {
	return models.ArangoID{}, errors.New("not implemented...")
}

func (session *TestSession) DropCollection(dbname string, collectionName string) error {
	_, ok := session.database[dbname]
	if !ok {
		return errors.New("DB doesnt")
	}
	delete(session.database[dbname], collectionName)
	return nil
}

func (session *TestSession) TruncateCollection(dbname string, collectionName string) error {
	session.database[dbname][collectionName] = make([]map[string]interface{}, 10)
	return nil
}

func (session *TestSession) CreateDocument(dbname string, collectionName string, object interface{}) (models.ArangoID, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	arangoID := models.ArangoID{
		ID:  timestamp,
		Key: strconv.FormatInt(rand.Int63(), 10),
		Rev: "",
	}

	// create entry
	entry := structs.Map(arangoID)
	for key, value := range objectToMap(object) {
		entry[key] = value
	}

	// "persist"
	session.database[dbname][collectionName] = append(session.database[dbname][collectionName], entry)

	return arangoID, nil
}

func (session *TestSession) CreateJSONDocument(dbname string, collectionName string, jsonObj string) (models.ArangoID, error) {
	jsonMap := make(map[string]interface{})
	json.Unmarshal([]byte(jsonObj), &jsonMap)
	return session.CreateDocument(dbname, collectionName, jsonMap)
}

func (session *TestSession) AqlQuery(dbname string, query string, count bool, batchSize int) ([]map[string]interface{}, string, error) {
	if len(session.aqlFakes) > 0 {
		aql := session.aqlFakes[query]
		return aql.MapResult, aql.JsonResult, nil
	}
	return nil, "nil", errors.New("fakes are empty")
}

func (session *TestSession) AddAqlFake(aql string, fake AqlFake) {
	session.aqlFakes[aql] = fake
}

func (session *TestSession) GetCollectionByID(dbname string, id string) (string, map[string]interface{}, error) {
	if entry := findByParam(session, dbname, "_id", id); entry != nil {
		jsonStr, err := json.Marshal(entry)
		return string(jsonStr), (*entry), err
	}
	return "", nil, errors.New("Cant find id")
}

func (session *TestSession) UpdateDocument(dbname string, id string, object interface{}) error {
	if entry := findByParam(session, dbname, "_id", id); entry != nil {
		for key, value := range objectToMap(object) {
			(*entry)[key] = value
		}
	}
	return nil
}

func (session *TestSession) UpdateJSONDocument(dbname string, id string, jsonObj string) error {
	jsonMap := make(map[string]interface{})
	json.Unmarshal([]byte(jsonObj), &jsonMap)
	return session.UpdateDocument(dbname, id, jsonMap)
}

func (tsession *TestSession) Migrate(migrations ...Migration) error {
	if list, _ := tsession.ListDBs(); !sliceTricks.Contains(list, migrationColl) {
		tsession.database[systemDB][migrationColl] = make([]map[string]interface{}, 20)
	}

	// helper function
	findMigration := func(name string) map[string]interface{} {
		if dbMigrations := tsession.database[systemDB][migrationColl]; len(systemDB) > 0 {
			for _, mig := range dbMigrations {
				if str, ok := mig["name"]; ok && str == name {
					return mig
				}
			}
		}
		return nil
	}

	// transform to session
	var session Session = tsession

	// iterate all migrations
	for _, mig := range migrations {
		if findMigration(mig.Name) == nil {
			mig.Handle(session)
			mig.Status = Finished
			tsession.database[systemDB][migrationColl] = append(tsession.database[systemDB][migrationColl], structs.Map(mig))
		}
	}
	return nil
}
