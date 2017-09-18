package aranGoConnection

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// AranGoConnection represent to connection
type AranGoConnection struct {
	Host      string
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

// HTTPDo function runs the HTTP request and processes its response in a new goroutine.
func (connection *AranGoConnection) HTTPDo(ctx context.Context, request *http.Request, processResponse func(*http.Response, error) error) error {

	// Run the HTTP request in a goroutine and pass the response to processResponse.
	transport := &http.Transport{}
	client := &http.Client{Transport: transport}
	errorChannel := make(chan error, 1)

	// do request
	go func() { errorChannel <- processResponse(client.Do(request)) }()
	select {
	case <-ctx.Done():
		transport.CancelRequest(request)
		<-errorChannel // wait for processResponse function
		return ctx.Err()
	case err := <-errorChannel:
		return err
	}
}

// Query queries the "route" with the http-"methode".
// The typ ist the pointer for your result
func (connection *AranGoConnection) Query(typ interface{}, methode string, route string, body interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// build request
	request, err := http.NewRequest(methode, connection.Host+route, bodyReader)
	if err != nil {
		return err
	}

	// set auth headers
	request.Header.Set("Content-Type", "application/json")
	if &connection.jwtString != nil {
		request.Header.Set("Authorization", "Bearer "+connection.jwtString)
	}

	// create http-Context
	httpContext, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFunc()

	// fire up request and unmarshal serverTime
	err = connection.HTTPDo(httpContext, request, func(response *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer response.Body.Close()
		decoder := json.NewDecoder(response.Body)
		if err := decoder.Decode(&typ); err != nil {
			return err
		}
		return nil
	})
	return err
}

// Post creates a POST-Request
func (connection *AranGoConnection) Post(url string, object interface{}) (string, map[string]interface{}, error) {
	// marshal body
	jsonBody, err := json.Marshal(object)
	failOnError(err, "Can't marshal object")
	return connection.PostJSON(url, jsonBody)
}

// PostJSON creates a POST-Request with JSONbody
func (connection *AranGoConnection) PostJSON(url string, jsonBody []byte) (string, map[string]interface{}, error) {
	// build url
	url = connection.Host + url

	// build request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", nil, err
	}
	return fireRequestAndUnmarshal(connection, req)
}

// Put creates a PUT-Request
func (connection *AranGoConnection) Put(url string, object interface{}) (string, map[string]interface{}, error) {
	// marshal body
	jsonBody, err := json.Marshal(object)
	failOnError(err, "Can't marshal object")
	return connection.PutJSON(url, jsonBody)
}

// PutJSON creates a PUT-Request with a jsonString
func (connection *AranGoConnection) PutJSON(url string, jsonBody []byte) (string, map[string]interface{}, error) {
	// build url
	url = connection.Host + url

	// build request
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", nil, err
	}
	return fireRequestAndUnmarshal(connection, req)
}

// Delete creates a DELETE-request
func (connection *AranGoConnection) Delete(url string) (string, map[string]interface{}, error) {
	// build url
	url = connection.Host + url

	// build request
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	failOnError(err, "Failed while building DELETE-request")

	return fireRequestAndUnmarshal(connection, req)
}

// Patch creates a PATCH-request
func (connection *AranGoConnection) Patch(url string, object interface{}) (string, map[string]interface{}, error) {
	// marshal body
	jsonBody, err := json.Marshal(object)
	failOnError(err, "Can't marshal object")
	return connection.PatchJSON(url, jsonBody)
}

// PatchJSON creates a PATCH-request with JSONbody
func (connection *AranGoConnection) PatchJSON(url string, jsonBody []byte) (string, map[string]interface{}, error) {
	// build url
	url = connection.Host + url

	// build request
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", nil, err
	}
	return fireRequestAndUnmarshal(connection, req)
}

func fireRequestAndUnmarshal(connection *AranGoConnection, request *http.Request) (string, map[string]interface{}, error) {
	// set headers
	request.Header.Set("Content-Type", "application/json")
	if &connection.jwtString != nil {
		request.Header.Set("Authorization", "Bearer "+connection.jwtString)
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// unmarshal to map
	var responseMap map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.UseNumber()
	err = decoder.Decode(&responseMap)

	return string(body), responseMap, err
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
