package aranGoDriver

import (
	"fmt"
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

// CreateDB test create a db
func (session *TestSession) CreateDB(dbname string) {
	session.database = append(session.database, dbname)
}
