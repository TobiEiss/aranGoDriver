package aranGoDriver

type Session interface {
	Connect(username string, password string) error

	// databases
	ListDBs() ([]string, error)
	CreateDB(dbname string) error
	DropDB(dbname string) error
	CreateCollection(dbname string, collectionName string) error
	DropCollection(dbname string, collectionName string) error
}
