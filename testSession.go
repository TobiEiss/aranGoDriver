package aranGoDriver

import (
	"fmt"
	"strconv"
	"time"

	"github.com/TobiEiss/aranGoDriver/models"

	"errors"
	"math/rand"

	"github.com/fatih/structs"
)

type TestSession struct {
	database map[string]map[string][]map[string]interface{}
}

func NewTestSession() *TestSession {
	// database - collection - list of document (key, value)
	return &TestSession{make(map[string]map[string][]map[string]interface{})}
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

func (session *TestSession) CreateDocument(dbname string, collectionName string, object map[string]interface{}) (models.ArangoID, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	arangoID := models.ArangoID{
		ID:  timestamp,
		Key: strconv.FormatInt(rand.Int63(), 10),
		Rev: "",
	}

	// create entry
	entry := structs.Map(arangoID)
	for key, value := range object {
		entry[key] = value
	}

	// "persist"
	session.database[dbname][collectionName] = append(session.database[dbname][collectionName], entry)

	return arangoID, nil
}

func (session *TestSession) AqlQuery(dbname string, query string, count bool, batchSize int) (map[string]interface{}, error) {
	// TODO
	return nil, nil
}
