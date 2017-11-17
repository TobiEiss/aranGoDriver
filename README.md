# aranGoDriver [![Build Status](https://travis-ci.org/TobiEiss/aranGoDriver.svg?branch=master)](https://travis-ci.org/TobiEiss/aranGoDriver)

This project is a golang-driver for [ArangoDB](https://www.arangodb.com/) writen in go.   
There is also an embedded-in-memory-Database to run all your tests.

Currently implemented:
* [version](#version)
* [databases](#database): connect, list, create, drop
* [collections](#collection): create, drop, truncate, update
* [documents](#document): create, getById
* [migrations](#migrations)
* [AQL](#aql): simple cursor

If you miss something, please contact me!

## Getting started
All you need is a running Arango-DB and a go-environment.

Install aranGoDriver:
`go get github.com/TobiEiss/aranGoDriver`

Write your first aranGoDriver-Programm:
```golang
func main() {
    var session aranGoDriver.Session
    
    // Initialize a arango-Session with the address to your arangoDB.
    //
    // If you write a test use:
    // session = aranGoDriver.NewTestSession()
    //
    session = aranGoDriver.NewAranGoDriverSession("http://localhost:8529")

    // Connect to your arango-database:
	session.Connect("usnername", "secretPassword")

    // Concrats, you are connected!
    // Let's print out all your databases
    list, err := session.ListDBs()
    if err != nil {
        log.Fatal("there was a problem: ", err)
    }
    log.Println(list)

    // Create a new database
    err = session.CreateDB("myNewDatabase")
    // TODO: handle err

    // Create a new collection
    err = session.CreateCollection("myNewDatabase", "myNewCollection")
    // TODO: handle err

    // Create Document
    newDocument := make(map[string]interface{})
    newDocument["foo"] = "bar"
    arangoID, err = session.CreateDocument("myNewDatabase", "myNewCollection", newDocument)
}
```

## Test

### Test against a fake-in-memory-database:
```
go test
```

### Test with a real database
```
go test -dbhost http://localhost:8529 -dbusername root -dbpassword password123
```

## Usage

### Connect to your ArangoDB

You need a new Session to your database with the hostname as parameter. Then connect with an existing username and a password.
```
session := aranGoDriver.NewAranGoDriverSession("http://localhost:8529")
session.Connect("username", "password")
```

### Version
```golang
version, err := session.Version()
```

### Database
```golang
// list databases
list := session.ListDBs()
fmt.Println(list) // will print ([]string): [ _system test testDB]

// create databases
err := session.CreateDB("myNewDatabase")

// drop databases
err = session.DropDB("myNewDatabase")
```

### Collection
```golang
// create a collection in a database
err = CreateCollection("myNewDatabase", "myNewCollection")

// drop collection from database
err = DropCollection("myNewDatabase", "myNewCollection")

// truncate database
err = TruncateCollection("myNewDatabase", "myNewCollection")
```

EdgeCollection:
```golang
// create a collection in a database
err = CreateEdgeCollection("myNewDatabase", "myNewEdgeCollection")
```

### Document
```golang
// create document
testDoc["foo"] = "bar"
arangoID, err := session.CreateDocument("myNewDatabase", "myNewCollection", testDoc)

// get by id
resultAsMap, err := session.GetCollectionByID("myNewDatabase", idOfDocument)

// update Document
testDoc["bar"] = "foo"
err = session.UpdateDocument("myNewDatabase", arangoID.ID, testDoc)
```

EdgeDocument
```golang
arangoID, err := session.CreateDocument("myNewDatabase", "myNewCollection", testDoc)
```

### Graphs
To create a graph, you need to define the edges and nodes for the graph.
This can be done with the `EdgeDefinition` model.
```golang
edgeDefinition := models.EdgeDefinition{
    Collection: "myEdgeCollection",
    From:       []string{"myCollection1"},
    To:         []string{"myCollection2"}}
edgeDefinitions := []models.EdgeDefinition{edgeDefinition}

err := session.CreateGraph("myDatabase", "myGraph", edgeDefinitions)
```

If you want to get rid of an existing graph, you can use the `DropGraph` method.
```golang
err := session.DropGraph("myDatabase", "myGraph")
```

For an overview of your existing graphs, you can use `ListGraphs`.
```golang
str, b, err := session.ListGraphs("myDatabase")
```

### Migrations
In some cases you need 'migrations'. For example, you need default-user in your database in every environment.
For this case, you can use migrations. The aranGoDriver write his own memo in a `migrations`-Collection in the standard `-system`-Database of arango, and execute the migration only one time.
AranGoDriver will identificate the migration by name.
Take a look to the following example:
```golang
// check migrations
mig1 := aranGoDriver.Migration{
    Name: "mig1", // name of migration to identificate
    Handle: func(embeddedSession aranGoDriver.Session) {
        testMap := make(map[string]interface{})
        testMap["foo"] = "foo"
        // you can do everything with the database
        embeddedSession.CreateDocument("myDatabase", "myCollection", testMap)
    },
}
// excute migration
// Run This line before you 'main-loop' of your program
session.Migrate(mig1)
```

### aql
```golang
// create query
query := "FOR element in testColl FILTER element.foo == 'bar' RETURN element"
response, err := session.AqlQuery("myNewDatabase", query, true, 1)
```
