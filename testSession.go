package aranGoDriver

import (
	"fmt"

	"github.com/TobiEiss/aranGoDriver/sliceTricks"
)

type TestSession struct {
	database []string
}

func NewTestSession() *TestSession {
	return &TestSession{}
}

// Connect test
func (session TestSession) Connect(username string, password string) {
	fmt.Println("Connect to DB")
}

func (session *TestSession) ListDBs() []string {
	return session.database
}

// CreateDB test create a db
func (session *TestSession) CreateDB(dbname string) {
	session.database = append(session.database, dbname)
}

func (session *TestSession) DropDB(dbname string) {
	index := sliceTricks.Find(session.database, func(index int, value string) bool {
		return value == dbname
	})
	session.database = append(session.database[:index], session.database[index+1:]...)
}
