package aranGoDriver

import (
	"log"

	"github.com/TobiEiss/aranGoDriver/aranGoConnection"
	"github.com/TobiEiss/aranGoDriver/models"
)

// AranGoSession represent to Session
type AranGoSession struct {
	arangoCon *aranGoConnection.AranGoConnection
}

const urlAuth = "/_open/auth"
const urlDatabase = "/_api/database"

// NewAranGoDriverSession creates a new instance of a AranGoDriver-Session.
// Need a host (e.g. "http://localhost:8529/")
func NewAranGoDriverSession(host string) *AranGoSession {
	return &AranGoSession{aranGoConnection.NewAranGoConnection(host)}
}

// Connect to arangoDB
func (session *AranGoSession) Connect(username string, password string) {
	credentials := models.Credentials{}
	credentials.Username = username
	credentials.Password = password

	resp := session.arangoCon.Post(urlAuth, credentials)
	session.arangoCon.SetJwtKey(resp["jwt"].(string))
}

// ListDBs lists all db's
func (session *AranGoSession) ListDBs() []string {
	resp := session.arangoCon.Get(urlDatabase)
	result := resp["result"].([]interface{})

	dblist := make([]string, len(result))
	for _, value := range result {
		if str, ok := value.(string); ok {
			dblist = append(dblist, str)
		}
	}

	return dblist
}

// CreateDB creates a new db
func (session *AranGoSession) CreateDB(dbname string) {
	body := make(map[string]string)
	body["name"] = dbname

	session.arangoCon.Post(urlDatabase, body)
}

// DropDB drop a database
func (session *AranGoSession) DropDB(dbname string) {
	session.arangoCon.Delete(urlDatabase + "/" + dbname)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
