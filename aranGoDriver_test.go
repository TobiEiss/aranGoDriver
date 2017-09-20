package aranGoDriver_test

import (
	"log"
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
		testSession := aranGoDriver.NewTestSession()

		// fake
		testDoc := make([]interface{}, 1)
		testMap := make(map[string]interface{})
		testMap["foo"] = "bar"
		testMap["_id"] = "userid"
		testDoc[0] = testMap
		fake1 := aranGoDriver.AqlFake{
			MapResult: testDoc,
		}
		testSession.AddAqlFake("FOR element in testColl FILTER element.foo == 'bar' RETURN element", fake1)
		session = testSession
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

	// Create Document
	testMap := make(map[string]interface{})
	testMap["foo"] = "bar"
	if *testDbHost == "" {
		testMap["_id"] = "userid"
	}
	_, err = session.CreateDocument(*testDbName, *testCollName, testMap)
	failOnError(err, "failed while creating Document")

	// session.AqlQuery
	query := "FOR element in testColl FILTER element.foo == 'bar' RETURN element"
	results := []map[string]interface{}{}
	err = session.AqlQuery(&results, *testDbName, query, true, 1)
	failOnError(err, "AQL-Query")
	assertTrue(len(results) > 0)
	t.Log(results)

	// search for ID
	id, ok := results[0]["_id"].(string)
	if !ok {
		t.Error("Cant find a key in result")
	}
	result, err := session.GetCollectionByID(*testDbName, id)
	id2, ok := result["_id"].(string)

	if !ok || (id != id2) {
		t.Error("id's arent the same")
	}

	// Update document
	result["bar"] = "foo"
	err = session.UpdateDocument(*testDbName, result["_id"].(string), result)
	failOnError(err, "failed while update document")
	// Verify
	result, err = session.GetCollectionByID(*testDbName, result["_id"].(string))
	assertTrue(result["bar"] == "foo")

	// check migrations
	mig1 := aranGoDriver.Migration{
		Name: "mig1",
		Handle: func(embeddedSession aranGoDriver.Session) {
			testMap := make(map[string]interface{})
			testMap["foo"] = "foo"
			embeddedSession.CreateDocument(*testDbName, *testCollName, testMap)
		},
	}
	// excute migration
	session.Migrate(mig1)

	// Drop collection
	err = session.DropCollection(*testDbName, *testCollName)
	failOnError(err, "Cant drop collection")

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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
