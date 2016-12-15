package aranGoDriver

type Session interface {
	Connect(username string, password string) error

	// databases
	ListDBs() ([]string, error)
	CreateDB(dbname string) error
	DropDB(dbname string) error
}
