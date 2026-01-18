package test

import (
	"database/sql"
	"log"
)

func TruncateUserProfile(dbConn *sql.DB) {
	truncateTable(dbConn, "user_profile")
}

func TruncateUserInvitation(dbConn *sql.DB) {
	truncateTable(dbConn, "user_invitation")
}

func truncateTable(dbConn *sql.DB, tableName string) {
	_, err := dbConn.Exec("TRUNCATE TABLE " + tableName + " CASCADE;")
	if err != nil {
		log.Printf("failed to truncate %s table: %v", tableName, err)
		return
	}
}
