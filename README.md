# aranGoDriver

This project is a golang-driver for [ArangoDB](https://www.arangodb.com/)

Currently implemented:
* connect to DB
* databases: list, create, drop

## TOC
- [Test](#test)
    - [Test against a fake-in-memory-database:](#test-against-a-fake-in-memory-database)
    - [Test with a real database](#test-with-a-real-database)
        - [fit tests with a real database](#fit-tests-with-a-real-database)
- [Usage](#usage)
    - [Connect to your ArangoDB](#connect-to-your-arangodb)
    - [List all database](#list-all-database)
    - [Create a new database](#create-a-new-database)
    - [Drop a database](#drop-a-database)


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

### List all database
You will get all databases as string-slice (`[]string`)
```
list := session.ListDBs()
fmt.Println(list)
```

will print:
`[ _system test testDB]`

### Create a new database
Just name it!
```
session.CreateDB("myNewDatabase")
```

### Drop a database
And now lets drop this database
```
session.DropDB("myNewDatabase")
```
