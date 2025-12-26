//go:build integration

package integration

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var code int
	defer func() {
		log.Printf("clean up complete")
		os.Exit(code)
	}()

	db, err := setupTestDB()
	if err != nil {
		log.Panicf("failed to start test container: %v", err)
	}
	defer teardownTestDB(db)

	setupTestDBConn(db)
	defer closeTestDBConn()

	dbConn := getTestDBConn()
	err = dbConn.Ping()
	if err != nil {
		log.Panicf("failed to ping testcontainer database: %v", err)
	}

	err = runMigrations(dbConn)
	if err != nil {
		log.Panicf("failed to run migrations: %v", err)
	}

	code = m.Run()
}
