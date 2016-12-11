package aranGoDriver

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"encoding/json"

	"github.com/TobiEiss/aranGoDriver/models"
)

// AranGoSession represent to Session
type AranGoSession struct {
	urlRoot   string
	jwtString string
}

const urlAuth = "/_open/auth"
const urlDatabase = "/_api/database"

// NewAranGoDriverSession creates a new instance of a AranGoDriver-Session.
// Need a host (e.g. "http://localhost:8529/")
func NewAranGoDriverSession(host string) *AranGoSession {
	return &AranGoSession{host, ""}
}

// Connect to arangoDB
func (session *AranGoSession) Connect(username string, password string) {
	credentials := models.Credentials{}
	credentials.Username = username
	credentials.Password = password

	resp := post(session, urlAuth, credentials)
	session.jwtString = resp["jwt"].(string)
	fmt.Println(session.jwtString)
}

func (session *AranGoSession) ListDBs() []string {
	resp := get(session, urlDatabase)
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

	resp := post(session, urlDatabase, body)
	fmt.Println(resp)
}

func get(session *AranGoSession, url string) map[string]interface{} {
	url = session.urlRoot + url
	fmt.Println("URL:>", url)

	// build request
	req, err := http.NewRequest("GET", url, nil)
	// use JWT-token if set
	if &session.jwtString != nil {
		req.Header.Set("Authorization", "Bearer "+session.jwtString)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println(req)

	failOnError(err, "Cant do post-request to "+url)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// unmarshal to map
	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	return responseMap
}

func post(session *AranGoSession, url string, object interface{}) map[string]interface{} {
	// marshal body
	jsonBody, err := json.Marshal(object)
	failOnError(err, "Cant marshal object")

	// build url
	url = session.urlRoot + url
	fmt.Println("URL:>", url)

	// build request
	var jsonString = []byte(jsonBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")
	// use JWT-token if set
	if &session.jwtString != nil {
		req.Header.Set("Authorization", "Bearer "+session.jwtString)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println(req)

	failOnError(err, "Cant do post-request to "+url)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// unmarshal to map
	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	return responseMap
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
