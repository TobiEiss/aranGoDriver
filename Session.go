package aranGoDriver

type Session interface {
	Connect(username string, password string)
	CreateDB(dbname string)
}
