package aranGoDriver

import (
	"fmt"

	"errors"
)

type TestSession struct {
	database map[string]map[string][]string
}

func NewTestSession() *TestSession {
	return &TestSession{make(map[string]map[string][]string)}
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
	session.database[dbname] = map[string][]string{}
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
	session.database[dbname][collectionName] = append(session.database[dbname][collectionName], collectionName)
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
	session.database[dbname][collectionName] = []string{}
	return nil
}
