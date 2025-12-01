package sqldb

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const _currentMigrationVersion = 1

func RunMigration(db *sql.DB, migrationFileUrl *string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	sourceUrl := "file://./migration"
	if migrationFileUrl != nil {
		sourceUrl = *migrationFileUrl
	}

	m, err := migrate.NewWithDatabaseInstance(sourceUrl, "postgres", driver)
	if err != nil {
		return err
	}

	return m.Migrate(_currentMigrationVersion)
}
