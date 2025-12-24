# Big Backend

Backend services built with my current preferences.  
Mainly a repo for me to refer to things I forget over time.

## Dev Tools
### Task
https://taskfile.dev/ is used as a helper tool to run common dev commands.
```terminaloutput
➜  bigbackend git:(main) task
task: [default] task -l
task: Available tasks for this project:
* down:             shut down database and remove orphans
* up:               spin up database
* map:gen:          Generate converters
* map:help:         Get current go version
* mig:create:       Create migration file. Usage: task mig:create name=<name>
* mig:run:          Run migrations according to _currentMigrationVersion
* sqlgen:gen:       Run go-jet generator to create SQL builder code
```
`Taskfile.yml` is the main task file.  
Sub task file can be found in `_taskfiles` folder.

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
It can also be run manually using `Taskfile.yml` commands, which is useful when for generating sql models and builders. 

### Database SQL Builder
Database SQL builder uses [go-jet/jet](https://github.com/go-jet/jet).  
Check out `Taskfile.yml` on how to generate models and sql builders.

## Mapper
[goverter](https://github.com/jmattheis/goverter) is used to generate mappers between different layers.  
Checkout `Taskfile.yml` on how to generate mappers.

## Plan
- [x] Setup logger
- [x] Generate model and sql builders
- [x] Custom struct with go-jet, ie: shop-spring decimal or date for example
- [x] Database get and create
- [x] Mapper
- [ ] Integrate db query, domain and handler, test and validation
  - [ ] creator
  - [ ] add deletedAt nullable
  - [ ] getter
  - [ ] updater
  - [ ] deleter
- [ ] Http server
  - [ ] Extract to standalone server instead of main
  - [ ] Switch to chi router
  - [ ] Add middleware, crash recovery, tracing
- [ ] Graceful shutdown (Done, but it's ugly improve it before marking as done)
- [ ] Replace mock config with real config, extracting from env vars
- [ ] More descriptive validator, self implement or maybe explore go-playground/validator
- [ ] Implement inbox and outbox pattern
- [ ] Increase test coverage
- [ ] Rename database
- [ ] Improve quality of `_dev/sqlgen/generator.go`
- [ ] Automate generation of Auditable methods on entities
