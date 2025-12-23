//go:build integration

package sqldb

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// SetupTestDB sets up a test Postgres database using testcontainers and returns a sql.DB connection.
// It also registers cleanup functions to terminate the container and close the database connection after the test.
func SetupTestDB(t *testing.T) *sql.DB {
	db, err := setupTestDB()
	if err != nil {
		log.Panicf("failed to start test container: %v", err)
	}
	t.Cleanup(func() {
		teardownTestDB(db)
	})

	connString, err := db.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		log.Panicf("failed to obtain connection string: %v", err)
	}
	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		log.Panicf("failed to connect to testcontainer database: %v", err)
	}
	t.Cleanup(func() {
		log.Printf("closing db connection")
		err := dbConn.Close()
		if err != nil {
			log.Printf("failed to close db connection: %v", err)
		}
	})

	err = dbConn.Ping()
	if err != nil {
		log.Panicf("failed to ping testcontainer database: %v", err)
	}

	err = RunMigration(dbConn, nil)
	if err != nil {
		log.Panicf("failed to run database migrations: %v", err)
	}

	return dbConn
}

func setupTestDB() (*postgres.PostgresContainer, error) {
	log.Println("setup test db")
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
