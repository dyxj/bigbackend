//go:build integration

package integration

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dyxj/bigbackend/pkg/sqldb"
	_ "github.com/lib/pq"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDbConn *sql.DB

// Test DB Container Management

func setupTestDB() (*postgres.PostgresContainer, error) {
	log.Println("setup test db container")
	ctx := context.Background()

	dbName := "bigbackend"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
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
		return postgresContainer, err
	}
	return postgresContainer, nil
}

func teardownTestDB(db *postgres.PostgresContainer) {
	log.Println("teardown test db")
	if err := testcontainers.TerminateContainer(db); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}

// Test DB Connection Management

func setupTestDBConn(db *postgres.PostgresContainer) {
	log.Println("setup test db connection")
	connString, err := db.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		log.Panicf("failed to obtain connection string: %v", err)
	}
	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		log.Panicf("failed to connect to testcontainer database: %v", err)
	}
	testDbConn = dbConn
}

func closeTestDBConn() {
	log.Println("close test db connection")
	if testDbConn != nil {
		err := testDbConn.Close()
		if err != nil {
			log.Printf("failed to close test db connection: %v", err)
		}
	}
}

// getTestDBConn returns global test database connection, closure is done via TestMain
func getTestDBConn() *sql.DB {
	return testDbConn
}

// Migrations

// runMigrations runs database migrations on the provided database connection.
func runMigrations(dbConn *sql.DB) error {
	projectRoot, err := getProjectRoot()
	if err != nil {
		return err
	}
	migrationPath := filepath.Join(projectRoot, "migration")
	migrationURL := "file://" + migrationPath

	err = sqldb.RunMigration(dbConn, &migrationURL)
	if err != nil {
		return err
	}
	return nil
}

// getProjectRoot returns the root directory of the project by locating the go.mod file.
func getProjectRoot() (string, error) {
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
