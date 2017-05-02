package framed

import "github.com/TobiEiss/aranGoDriver"

type FramedConnection struct {
	Session *aranGoDriver.Session
}

// NewFramedConnection creates a new framedConnection which wraps a session.
func NewFramedConnection(session aranGoDriver.Session) *FramedConnection {
	return &FramedConnection{Session: &session}
}

// DB returns a database-model.
func (connection *FramedConnection) DB(database string) *Database {
	return &Database{Name: database, Session: connection.Session}
}

// CreateDB creates a new database
func (connection *FramedConnection) CreateDB(database string) (*Database, error) {
	err := (*connection.Session).CreateDB(database)
	return connection.DB(database), err
}

// DropDB drops the database
func (connection *FramedConnection) DropDB(database *Database) error {
	return (*connection.Session).DropDB(database.Name)
}
