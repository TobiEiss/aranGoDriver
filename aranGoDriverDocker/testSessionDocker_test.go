package aranGoDriverDocker_test

import (
	"log"
	"testing"

	"github.com/TobiEiss/aranGoDriver/aranGoDriverDocker"
)

func TestArangoDBDocker(t *testing.T) {
	password := "password"
	session, closeFunc := aranGoDriverDocker.SetupDockerTest(password)
	defer closeFunc()

	_, err := (*session).ListDBs()
	if err != nil {
		log.Println(err)
		t.Fail()
	}
}

func TestArangoDBDockerVersion(t *testing.T) {
	password := "password"
	session, closeFunc := aranGoDriverDocker.SetupDockerTest(password)
	defer closeFunc()

	// check Version
	version, err := (*session).Version()
	if version.Server == "" || err != nil {
		log.Println(err)
		t.Fail()
	}
}
