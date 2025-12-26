//go:build integration

package testx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Environment represents a test environment with optional database and server components.
// In the event there is a need to run integration test in parallel there is an option to spin up multiple environments
type Environment struct {
	name         string
	enableServer bool

	dbContainer *postgres.PostgresContainer
	dbConn      *sql.DB

	ready     chan struct{}
	errorChan chan error
	close     chan struct{}
	wgClose   *sync.WaitGroup

	runOnce   sync.Once
	closeOnce sync.Once

	logger *log.Logger
}

func NewEnvironment(name string, enableServer bool) *Environment {
	return &Environment{
		name:         name,
		enableServer: enableServer,
		ready:        make(chan struct{}),
		errorChan:    make(chan error),
		close:        make(chan struct{}),
		wgClose:      &sync.WaitGroup{},
		logger:       log.New(os.Stderr, fmt.Sprintf("test-env-%s ", name), log.LstdFlags),
	}
}

func (e *Environment) DBConn() *sql.DB {
	return e.dbConn
}

// Run starts the environment and returns channels to signal readiness and errors.
// If it is already running, it simply returns the existing channels.
func (e *Environment) Run() (<-chan struct{}, <-chan error) {
	e.runOnce.Do(func() {
		e.logger.Printf("starting environment")
		go e.run(e.ready, e.errorChan)
	})

	return e.ready, e.errorChan
}

func (e *Environment) run(ready chan struct{}, errorChan chan error) {
	err := e.setupDBContainer()
	if err != nil {
		errorChan <- err
		return
	}
	e.wgClose.Add(1)
	defer e.teardownDBContainer()

	err = e.setupDBConn()
	if err != nil {
		errorChan <- err
		return
	}
	e.wgClose.Add(1)
	defer e.closeDBConn()

	err = e.runMigrations()
	if err != nil {
		errorChan <- err
		return
	}

	close(ready)
	<-e.close
}

// Close signal environment to close resources and waits for cleanup to complete.
func (e *Environment) Close() {
	e.logger.Printf("closing environment")
	e.closeOnce.Do(func() {
		close(e.close)
	})
	e.wgClose.Wait()
	e.logger.Printf("environment closed")
}

// setupDBContainer sets up the database container.
func (e *Environment) setupDBContainer() error {
	e.logger.Printf("setup db container")

	ctx := context.Background()

	dbName := "bigbackend"
	dbUser := "user"
	dbPassword := "password"

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return err
	}
	e.dbContainer = pgContainer
	return nil
}

// teardownDBContainer terminates the database container.
func (e *Environment) teardownDBContainer() {
	e.logger.Printf("teardown db container")
	if err := testcontainers.TerminateContainer(e.dbContainer); err != nil {
		e.logger.Printf("failed to teardown container: %v", err)
	}
	e.wgClose.Done()
}

// setupDBConn sets up the database connection using the connection string from the dbContainer.
func (e *Environment) setupDBConn() error {
	e.logger.Printf("setup db connection")
	connString, err := e.dbContainer.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		return err
	}
	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}
	e.dbConn = dbConn
	return nil
}

// closeDBConn closes the database connection if it is not nil.
func (e *Environment) closeDBConn() {
	e.logger.Printf("close db connection")
	err := e.dbConn.Close()
	if err != nil {
		e.logger.Printf("failed to close db connection: %v", err)
	}
	e.wgClose.Done()
}

// runMigrations runs database migrations on the provided database connection.
func (e *Environment) runMigrations() error {
	projectRoot, err := e.getProjectRoot()
	if err != nil {
		return err
	}
	migrationPath := filepath.Join(projectRoot, "migration")
	migrationURL := "file://" + migrationPath

	err = sqldb.RunMigration(e.dbConn, &migrationURL)
	if err != nil {
		return err
	}
	return nil
}

// getProjectRoot returns the root directory of the project by locating the go.mod file.
func (e *Environment) getProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("go.mod not found in any parent directory")
		}
		dir = parent
	}
}
