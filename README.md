# Big Backend

Backend services built with my current preferences.  
Mainly a repo for me to refer to things I forget over time.  

## Notes
It should be noted, not all projects require such levels of abstractions and layers, just covering all scenarios.  
The intent is also to cover microservices scenarios, however written in repo for ease of reference.

## Dev Tools
- Docker and docker compose, should come together with Docker Desktop
- [Taskfile](https://taskfile.dev/)
  - As it is heavily used for development, installing it with autocompletion is recommended.

### Task
`Taskfile.yml` is the main task file.  
Sub task file can be found in `_taskfiles` folder.
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
It can also be run manually using `Taskfile.yml` commands, which is useful when generating sql models and builders.  

Automatic migration should be used with caution and consider scenarios where multiple pods are configured. Ensure 
appropriate techniques are employed to handle these scenarios.

### Database SQL Builder
Database SQL builder uses [go-jet/jet](https://github.com/go-jet/jet).  
Check out `Taskfile.yml` on how to generate models and sql builders.

## Mapper
[goverter](https://github.com/jmattheis/goverter) is used to generate mappers between different layers.  
Checkout `Taskfile.yml` on how to generate mappers.  
Examples can be found in `[domain]_mapdef.go` files, resulting generated file is `[domain]_mapper.go`.

## Next Steps
- [ ] Replace mock config with real config, extracting from env vars
- [ ] Http server
  - [x] Extract to standalone server instead of main
  - [ ] Switch to chi router
  - [ ] Add middleware, crash recovery, tracing
- [ ] Implement inbox and outbox pattern
  - [ ] Consider watermill
- [ ] Scheduled job orchestration
- [ ] More descriptive validator, self implement or maybe explore go-playground/validator
- [ ] Improve quality of `_dev/sqlgen/generator.go`
- [ ] Automate generation of Auditable methods on entities
- [ ] Taskfile improvements
  - [ ] Move DB up and down sub command
  - [ ] Add run commands
  - [ ] Add test commands customized with build tags
    - not as urgent as currently use Goland to run
- [ ] Diagrams
  - Revise plantUML, completely forgot syntax
- [ ] updater
- [ ] deleter, add deletedAt nullable and indexing
- [ ] Improve error structure and packages, there are properties that do the same things
  - ie http status code and error response code
  - error are also split in multiple packages, can be done better