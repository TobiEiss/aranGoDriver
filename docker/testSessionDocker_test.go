package aranGoDriver_test

import (
	"log"
	"testing"

	aranGoDriver "github.com/TobiEiss/aranGoDriver/docker"
)

func TestArangoDBDocker(t *testing.T) {
	password := "password"
	session, closeFunc := aranGoDriver.SetupDockerTest(password)
	//time.Sleep(time.Microsecond * 10000)
	defer closeFunc()

	_, err := (*session).ListDBs()
	if err != nil {
		log.Println(err)
		t.Fail()
	}
}
