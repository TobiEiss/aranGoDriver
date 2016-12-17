package aranGoConnection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// AranGoConnection represent to connection
type AranGoConnection struct {
	urlRoot   string
	jwtString string
}

// NewAranGoConnection creates a new instance of a AranGoDriver-Connection.
// Need a host (e.g. "http://localhost:8529/")
func NewAranGoConnection(host string) *AranGoConnection {
	return &AranGoConnection{host, ""}
}

// SetJwtKey sets the JWT-Token
func (connection *AranGoConnection) SetJwtKey(jwtString string) {
	connection.jwtString = jwtString
}

// Get creates a GET-Request
func (connection *AranGoConnection) Get(url string) (map[string]interface{}, error) {
	url = connection.urlRoot + url
	fmt.Println("URL:>", url)

	// build request
	req, err := http.NewRequest("GET", url, nil)
	failOnError(err, "Failed while build GET-request")

	return fireRequestAndUnmarshal(connection, req)
}

// Post creates a POST-Request
func (connection *AranGoConnection) Post(url string, object interface{}) (map[string]interface{}, error) {
	// marshal body
	jsonBody, err := json.Marshal(object)
	failOnError(err, "Cant marshal object")

	// build url
	url = connection.urlRoot + url
	fmt.Println("URL:>", url)

	// build request
	var jsonString = []byte(jsonBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	return fireRequestAndUnmarshal(connection, req)
}

// Put creates a PUT-Request
func (connection *AranGoConnection) Put(url string, object interface{}) (map[string]interface{}, error) {
	// marshal body
	jsonBody, err := json.Marshal(object)
	failOnError(err, "Cant marshal object")

	// build url
	url = connection.urlRoot + url
	fmt.Println("URL:>", url)

	// build request
	var jsonString = []byte(jsonBody)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonString))
	return fireRequestAndUnmarshal(connection, req)
}

// Delete creates a DELETE-request
func (connection *AranGoConnection) Delete(url string) (map[string]interface{}, error) {
	// build url
	url = connection.urlRoot + url
	fmt.Println("URL:>", url)

	// build request
	req, err := http.NewRequest("DELETE", url, nil)
	failOnError(err, "Failed while build DELETE-request")

	return fireRequestAndUnmarshal(connection, req)
}

func fireRequestAndUnmarshal(connection *AranGoConnection, request *http.Request) (map[string]interface{}, error) {
	// set headers
	request.Header.Set("Content-Type", "application/json")
	if &connection.jwtString != nil {
		request.Header.Set("Authorization", "Bearer "+connection.jwtString)
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// unmarshal to map
	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)

	return responseMap, err
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
