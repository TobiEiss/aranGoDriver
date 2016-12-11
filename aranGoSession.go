package aranGoDriver

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type AranGoSession struct {
	urlRoot   string
	jwtString string
}

const urlAuth = "/_open/auth"

func NewAranGoDriverSession(host string) *AranGoSession {
	return &AranGoSession{host, ""}
}

func (session AranGoSession) Connect(username string, password string) {
	resp := post(&session, urlAuth, "{\"username\":\""+username+"\", \"password\":\""+password+"\"}")
	fmt.Println(resp)
}

func post(session *AranGoSession, url string, jsonBody string) string {
	url = session.urlRoot + url
	fmt.Println("URL:>", url)

	var jsonString = []byte(jsonBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	failOnError(err, "Cant do post-request to "+url)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
