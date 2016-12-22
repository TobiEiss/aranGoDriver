package aranGoDriver

import (
	"github.com/TobiEiss/aranGoDriver/models"
)

// MigrationStatus is like a enum and represent the status
type MigrationStatus string

const (
	// Started means, the migration has started
	Started MigrationStatus = "started"
	// Finished means, the migration has finished
	Finished MigrationStatus = "finished"
)

// MigrationExecute is the func to execute the migration
type MigrationExecute func(Session)

// Migration represent the whole migration
type Migration struct {
	models.ArangoID `json:"-"`
	Name            string           `json:"name"`
	Handle          MigrationExecute `json:"-"`
	Status          MigrationStatus  `json:"status"`
}
