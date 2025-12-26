package integration

import (
	"database/sql"
	"log"
)

func truncateUserProfile(dbConn *sql.DB) {
	log.Printf("truncating user_profile table")
	_, err := dbConn.Exec("TRUNCATE TABLE user_profile CASCADE;")
	if err != nil {
		log.Printf("failed to truncate user_profile table: %v", err)
		return
	}
}
