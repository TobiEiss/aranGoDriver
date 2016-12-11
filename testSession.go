package aranGoDriver

type TestSession struct {
}

func NewTestSession() *TestSession {
	return &TestSession{}
}

func (session TestSession) Connect(username string, password string) {

}

func (session TestSession) CreateDB(dbname string) {

}
