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
	return &Database{Name: database}
}
