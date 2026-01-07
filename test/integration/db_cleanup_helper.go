package integration

import (
	"database/sql"
	"log"
)

func truncateUserProfile(dbConn *sql.DB) {
	truncateTable(dbConn, "user_profile")
}

func truncateUserInvitation(dbConn *sql.DB) {
	truncateTable(dbConn, "user_invitation")
}

func truncateTable(dbConn *sql.DB, tableName string) {
	log.Printf("truncating %s table", tableName)
	_, err := dbConn.Exec("TRUNCATE TABLE " + tableName + " CASCADE;")
	if err != nil {
		log.Printf("failed to truncate %s table: %v", tableName, err)
		return
	}
}
