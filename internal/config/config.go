package config

import "time"

type Config struct {
	HTTPServerConfig *HTTPServerConfig
	DBConfig         *DBConfig
}

type HTTPServerConfig struct {
	isDebug             bool
	host                string
	port                int
	readHeaderTimeout   time.Duration
	readTimeout         time.Duration
	idleTimeout         time.Duration
	handlerTimeout      time.Duration
	shutdownTimeout     time.Duration
	shutdownHardTimeout time.Duration
	readinessProbeDelay time.Duration
}

func (c *HTTPServerConfig) IsDebug() bool {
	return c.isDebug
}

func (c *HTTPServerConfig) Host() string {
	return c.host
}

func (c *HTTPServerConfig) Port() int {
	return c.port
}

func (c *HTTPServerConfig) ReadHeaderTimeout() time.Duration {
	return c.readHeaderTimeout
}

func (c *HTTPServerConfig) ReadTimeout() time.Duration {
	return c.readTimeout
}

func (c *HTTPServerConfig) IdleTimeout() time.Duration {
	return c.idleTimeout
}

func (c *HTTPServerConfig) HandlerTimeout() time.Duration {
	return c.handlerTimeout
}

func (c *HTTPServerConfig) ShutDownTimeout() time.Duration {
	return c.shutdownTimeout
}

func (c *HTTPServerConfig) ShutDownHardTimeout() time.Duration {
	return c.shutdownHardTimeout
}

func (c *HTTPServerConfig) ReadinessProbeDelay() time.Duration {
	return c.readinessProbeDelay
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
		isDebug:             false,
		host:                "",
		port:                8080,
		readHeaderTimeout:   500 * time.Millisecond,
		readTimeout:         500 * time.Millisecond,
		idleTimeout:         time.Second,
		handlerTimeout:      2 * time.Second,
		shutdownTimeout:     15 * time.Second,
		shutdownHardTimeout: 3 * time.Second,
		readinessProbeDelay: 5 * time.Second,
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
