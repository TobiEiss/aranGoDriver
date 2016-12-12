package aranGoDriver_test

import (
	"testing"

	"flag"

	"github.com/TobiEiss/aranGoDriver"
	"github.com/TobiEiss/aranGoDriver/sliceTricks"
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

	// Connect
	session.Connect("root", "ILoBhREd36LB8USwpcHcCz4hLjj8k")

	// Check listDB
	list := session.ListDBs()
	assertTrue(!sliceTricks.Contains(list, "testDB"))
	t.Log(list)

	// CreateDB
	session.CreateDB("testDB")
	list = session.ListDBs()
	t.Log(session.ListDBs())
	assertTrue(sliceTricks.Contains(list, "testDB"))

	// DropDB
	session.DropDB("testDB")
	list = session.ListDBs()
	t.Log(session.ListDBs())
	assertTrue(!sliceTricks.Contains(list, "testDB"))
}

func assertTrue(test bool) {
	if !test {
		panic("Assertion failed")
	}
}
