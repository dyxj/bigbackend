package config

type DBConfig struct {
	HostEV     string `env:"DB_HOST"`
	PortEV     int    `env:"DB_PORT"`
	UserEV     string `env:"DB_USER"`
	PasswordEV string `env:"DB_PASSWORD"`
	DBNameEV   string `env:"DB_NAME"`
}

func (c *DBConfig) Host() string {
	return c.HostEV
}

func (c *DBConfig) Port() int {
	return c.PortEV
}

func (c *DBConfig) User() string {
	return c.UserEV
}

func (c *DBConfig) Password() string {
	return c.PasswordEV
}

func (c *DBConfig) DBName() string {
	return c.DBNameEV
}
