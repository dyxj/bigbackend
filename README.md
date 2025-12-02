# Big Backend

Backend services built with my current preferences.  
Mainly a repo for me to refer to things I forget over time.

## Database
PostgreSQL is used as the main database. Check out `Taskfile.yml` on how to spin up local database.

### Migration scripts
Database migration uses [golang-migrate/migrate](https://github.com/golang-migrate/migrate).

#### Creating migration scripts  
Have a look at `_taskfiles/migrate.yml` for available commands, mainly how to generate new migration files.  
Ensure migrations are always wrapped with a transaction.  

#### Running migration scripts
Migration scripts are executed automatically on application start up.
`_currentMigrationVersion` in `pkg/sqldb/migrate.go` determines the migration scripts to run.

### Database SQL Builder

## Plan
- [x] Setup logger
- [] Database get and create
- [] Http server
- [] Replace mock config with real config, extracting from env vars
- [] Custom struct with go-jet, ie: shop-spring decimal
