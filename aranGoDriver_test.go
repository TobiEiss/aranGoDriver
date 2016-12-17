package aranGoDriver_test

import (
	"testing"

	"flag"

	"github.com/TobiEiss/aranGoDriver"
	"github.com/TobiEiss/aranGoDriver/sliceTricks"
)

var (
	testDbHost   = flag.String("dbhost", "", "run database integration test")
	testUsername = flag.String("dbusername", "testUser", "username of test-user")
	testPassword = flag.String("dbpassword", "password123", "password for test-user")
	testDbName   = flag.String("dbtestdbname", "testDB", "database name of test-database")
	testCollName = flag.String("dbtestcollname", "testColl", "collection name of test-collection")
)

func TestMain(t *testing.T) {
	flag.Parse()
	t.Log("Start tests")

	var session aranGoDriver.Session

	// check the flag database
	if *testDbHost != "" {
		t.Log("use arangoDriver")
		session = aranGoDriver.NewAranGoDriverSession(*testDbHost)
	} else {
		t.Log("use testDriver")
		session = aranGoDriver.NewTestSession()
	}

	// Connect
	session.Connect(*testUsername, *testPassword)

	// DropDB
	err := session.DropDB(*testDbName)
	t.Log(err)

	// Check listDB
	list, err := session.ListDBs()
	t.Log(err)
	assertTrue(!sliceTricks.Contains(list, *testDbName))
	t.Log(list)

	// CreateDB
	err = session.CreateDB(*testDbName)
	list, err = session.ListDBs()
	t.Log(session.ListDBs())
	assertTrue(sliceTricks.Contains(list, *testDbName))

	// DropDB
	err = session.DropDB(*testDbName)
	list, err = session.ListDBs()
	t.Log(session.ListDBs())
	assertTrue(!sliceTricks.Contains(list, *testDbName))

	// CreateDB
	err = session.CreateDB(*testDbName)
	list, err = session.ListDBs()
	t.Log(session.ListDBs())
	assertTrue(sliceTricks.Contains(list, *testDbName))

	// Create collection
	err = session.CreateCollection(*testDbName, *testCollName)

	// Truncate collection
	err = session.TruncateCollection(*testDbName, *testCollName)

	// Drop collection
	err = session.DropCollection(*testDbName, *testCollName)
	t.Log(err)

	// DropDB
	err = session.DropDB(*testDbName)
	list, err = session.ListDBs()
	t.Log(session.ListDBs())
	assertTrue(!sliceTricks.Contains(list, *testDbName))
}

func assertTrue(test bool) {
	if !test {
		panic("Assertion failed")
	}
}
