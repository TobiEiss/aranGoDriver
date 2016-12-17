# aranGoDriver

This project is a golang-driver for [ArangoDB](https://www.arangodb.com/)

Currently implemented:
* connect to DB
* databases: list, create, drop
* collections: create, drop, truncate

## TOC
- [Test](#test)
    - [Test against a fake-in-memory-database:](#test-against-a-fake-in-memory-database)
    - [Test with a real database](#test-with-a-real-database)
- [Usage](#usage)
    - [Connect to your ArangoDB](#connect-to-your-arangodb)
    - [Database](#database)
    - [Collection](#collection)

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
```
// list databases
list := session.ListDBs()
fmt.Println(list) // will print ([]string): [ _system test testDB]

// create databases
session.CreateDB("myNewDatabase")

// drop databases
session.DropDB("myNewDatabase")
```

### Collection
```
// create a collection in a database
CreateCollection("myNewDatabase", "myNewCollection")

// drop collection from database
DropCollection("myNewDatabase", "myNewCollection")

// truncate database
TruncateCollection("myNewDatabase", "myNewCollection")
```
