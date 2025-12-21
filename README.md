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
- [x] Generate model and sql builders
- [x] Custom struct with go-jet, ie: shop-spring decimal or date for example
- [x] Database get and create
- [x] Mapper
- [] Integrate db query, domain and handler, test and validation
- [] Http server
- [] Graceful shutdown (Done, but it's ugly improve it before marking as done)
- [] Replace mock config with real config, extracting from env vars
- [] Improve quality of sql gen generated code
- [] Tests
- [] Rename database
