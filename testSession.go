package aranGoDriver

import (
	"fmt"

	"errors"

	"github.com/TobiEiss/aranGoDriver/sliceTricks"
)

type TestSession struct {
	database []string
}

func NewTestSession() *TestSession {
	return &TestSession{}
}

// Connect test
func (session TestSession) Connect(username string, password string) error {
	fmt.Println("Connect to DB")
	return nil
}

func (session *TestSession) ListDBs() ([]string, error) {
	return session.database, nil
}

// CreateDB test create a db
func (session *TestSession) CreateDB(dbname string) error {
	if sliceTricks.Contains(session.database, dbname) {
		return errors.New("")
	}
	session.database = append(session.database, dbname)
	return nil
}

func (session *TestSession) DropDB(dbname string) error {
	index := sliceTricks.Find(session.database, func(index int, value string) bool {
		return value == dbname
	})
	if index < 0 {
		return errors.New("")
	}
	session.database = append(session.database[:index], session.database[index+1:]...)
	return nil
}
