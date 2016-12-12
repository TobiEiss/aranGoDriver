package aranGoDriver_test

import (
	"testing"

	"flag"

	"github.com/TobiEiss/aranGoDriver"
	"github.com/TobiEiss/aranGoDriver/sliceTricks"
)

const testUsername string = "root"
const testPassword string = "ILoBhREd36LB8USwpcHcCz4hLjj8k"
const testDbName string = "testDB"
const testDbHost string = "http://arangodb:8529"

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
		session = aranGoDriver.NewAranGoDriverSession(testDbHost)
	} else {
		t.Log("use testDriver")
		session = aranGoDriver.NewTestSession()
	}

	// Connect
	session.Connect(testUsername, testPassword)

	// Check listDB
	list := session.ListDBs()
	assertTrue(!sliceTricks.Contains(list, testDbName))
	t.Log(list)

	// CreateDB
	session.CreateDB(testDbName)
	list = session.ListDBs()
	t.Log(session.ListDBs())
	assertTrue(sliceTricks.Contains(list, testDbName))

	// DropDB
	session.DropDB(testDbName)
	list = session.ListDBs()
	t.Log(session.ListDBs())
	assertTrue(!sliceTricks.Contains(list, testDbName))
}

func assertTrue(test bool) {
	if !test {
		panic("Assertion failed")
	}
}
