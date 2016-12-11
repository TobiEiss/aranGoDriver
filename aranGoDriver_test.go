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
	list := session.ListDBs()
	assertTrue(!contains(list, "testDB"))
	t.Log(list)
	session.CreateDB("testDB")
	t.Log(session.ListDBs())
}

func assertTrue(test bool) {
	if !test {
		panic("Assertion failed")
	}
}

func contains(stringList []string, searchString string) bool {
	for _, value := range stringList {
		if value == searchString {
			return true
		}
	}
	return false
}
