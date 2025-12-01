package sqldb

type Config interface {
	Host() string
	Port() int
	User() string
	Password() string
	DBName() string
}
