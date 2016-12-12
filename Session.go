package aranGoDriver

type Session interface {
	Connect(username string, password string)

	// databases
	ListDBs() []string
	CreateDB(dbname string)
	DropDB(dbname string)
}
