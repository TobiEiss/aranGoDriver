package aranGoDriver

import (
	"log"
	"time"

	dockertest "gopkg.in/ory-am/dockertest.v3"
)

/*
SetupDockerTest setups a arangoDB-docker-image, connects a session and return this Session.
Also retuns a purge-function to defer.

password := "password"
session, closeFunc := aranGoDriver.SetupDockerTest(password)
defer closeFunc()

err := session.Connect("root", "password")
*/
func SetupDockerTest(arangoDbPassword string) (*Session, func()) {
	var session Session
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an arangodb image, creates a container based on it and runs it
	resource, err := pool.Run("arangodb", "3.2.2", []string{"ARANGO_ROOT_PASSWORD=" + arangoDbPassword})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		session = NewAranGoDriverSession("http://localhost:" + resource.GetPort("8529/tcp"))
		// TODO: find a better solution instead of "sleep"
		time.Sleep(time.Millisecond * 1000)
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return &session, func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}
