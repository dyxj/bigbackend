package config

type DBConfig struct {
	HostEV     string `env:"HOST"`
	PortEV     int    `env:"PORT"`
	UserEV     string `env:"USER"`
	PasswordEV string `env:"PASSWORD"`
	DBNameEV   string `env:"NAME"`
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
