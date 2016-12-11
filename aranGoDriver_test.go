package aranGoDriver_test

import (
	"testing"

	"flag"

	"github.com/TobiEiss/aranGoDriver"
)

var (
	database = flag.Bool("database", false, "run database integration test")
)

func TestMain(t *testing.T) {
	flag.Parse()
	t.Log("Start tests")

	var session aranGoDriver.Session

	// check the flag database
	if *database {
		t.Log("use arangoDriver")
		session = aranGoDriver.NewAranGoDriverSession("http://arangodb:8529")
	} else {
		t.Log("use testDriver")
		session = aranGoDriver.NewTestSession()
	}

	session.Connect("root", "ILoBhREd36LB8USwpcHcCz4hLjj8k")
	t.Log(session.ListDBs())
	session.CreateDB("testDB")
	t.Log(session.ListDBs())
}
