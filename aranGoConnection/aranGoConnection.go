package aranGoConnection

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
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
