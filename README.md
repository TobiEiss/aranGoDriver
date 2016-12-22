# aranGoDriver [![Build Status](https://travis-ci.org/TobiEiss/aranGoDriver.svg?branch=master)](https://travis-ci.org/TobiEiss/aranGoDriver)

This project is a golang-driver for [ArangoDB](https://www.arangodb.com/)

Currently implemented:
* connect to DB
* databases: list, create, drop
* collections: create, drop, truncate, update
* documents: create, getById
* AQL: simple cursor

## TOC
- [Test](#test)
    - [Test against a fake-in-memory-database:](#test-against-a-fake-in-memory-database)
    - [Test with a real database](#test-with-a-real-database)
- [Usage](#usage)
    - [Connect to your ArangoDB](#connect-to-your-arangodb)
    - [Database](#database)
    - [Collection](#collection)
    - [Document](#document)
    - [AQL](#aql)

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

### Document
```golang
// create document
testDoc["foo"] = "bar"
arangoID, err := session.CreateDocument("myNewDatabase", "myNewCollection", testDoc)

// get by id
resultAsJsonString, resultAsMap, err := session.GetCollectionByID("myNewDatabase", idOfDocument)

// update Document
testDoc["bar"] = "foo"
err = session.UpdateDocument("myNewDatabase", arangoID.ID, testDoc)
```

### aql
```golang
// create query
query := "FOR element in testColl FILTER element.foo == 'bar' RETURN element"
response, err := session.AqlQuery("myNewDatabase", query, true, 1)
```
