package sqldb

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const _currentMigrationVersion = 3

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

	err = m.Migrate(_currentMigrationVersion)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}

	return nil
}
