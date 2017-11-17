package aranGoDriverDocker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/TobiEiss/aranGoDriver"

	dockertest "gopkg.in/ory-am/dockertest.v3"
)

var wg sync.WaitGroup

/*
SetupDockerTest setups a arangoDB-docker-image, connects a session and return this Session.
Also retuns a purge-function to defer.

password := "password"
session, closeFunc := aranGoDriver.SetupDockerTest(password)
defer closeFunc()
*/
func SetupDockerTest(arangoDbPassword string) (*aranGoDriver.Session, func()) {
	var session aranGoDriver.Session
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
		session = aranGoDriver.NewAranGoDriverSession("http://localhost:" + resource.GetPort("8529/tcp"))

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		wg.Add(1)
		go waitForConnection(ctx, func() error { return session.Connect("root", arangoDbPassword) })
		wg.Wait()

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

func waitForConnection(ctx context.Context, connectFunc func() error) error {
	defer wg.Done()

	for {
		select {
		case <-time.After(1 * time.Second):
			err := connectFunc()
			if err == nil {
				return nil
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
