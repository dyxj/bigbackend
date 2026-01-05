package main

import (
	"context"
	"log"

	"github.com/dyxj/bigbackend/internal/config"
	"github.com/dyxj/bigbackend/pkg/sqldb"
)

func main() {
	// Parse environment variables
	cfg, err := config.LoadDBConfig()
	if err != nil {
		log.Panicf("failed to load config: %v", err)
	}

	// Initialize database connection
	dbConn, err := sqldb.NewDBConn(context.Background(), cfg)
	if err != nil {
		log.Panicf("failed to connect to database: %v", err)
	}
	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Printf("failed to close db conn: %v", err)
		}
	}()

	// Run database migrations
	err = sqldb.RunMigration(dbConn, nil)
	if err != nil {
		log.Panicf("failed to run database migrations: %v", err)
	}
}
