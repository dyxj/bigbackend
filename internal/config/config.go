package config

type Config struct {
	HTTPServerConfig *HTTPServerConfig
	DBConfig         *DBConfig
}

type HTTPServerConfig struct {
	host string
	port int
}

type DBConfig struct {
	host     string
	port     int
	user     string
	password string
	dbName   string
}

func (c *DBConfig) Host() string {
	return c.host
}

func (c *DBConfig) Port() int {
	return c.port
}

func (c *DBConfig) User() string {
	return c.user
}

func (c *DBConfig) Password() string {
	return c.password
}

func (c *DBConfig) DBName() string {
	return c.dbName
}

func LoadConfig() (*Config, error) {
	httpConfig := &HTTPServerConfig{
		host: "",
		port: 8080,
	}
	dbConfig := &DBConfig{
		host:     "localhost",
		port:     5430,
		user:     "bigbackend_role",
		password: "postgrespw",
		dbName:   "bigbackend",
	}
	return &Config{
		HTTPServerConfig: httpConfig,
		DBConfig:         dbConfig,
	}, nil
}
