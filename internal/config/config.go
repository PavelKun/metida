package config

import (
	"fmt"
)

//Todo: реализовать возможность наполнение конфига из файла конфигурации через флаг
type dataBaseConnection struct {
	Host     string
	Port     int
	Table    string
	User     string
	Password string
}

// Config contains all necessary params.
type ConfigI interface {
	GetCoonnect() string
	CommandLineI
}

// Config contains all necessary params.
type Config struct {
	db dataBaseConnection
	CommandLineI
}

// NewConfig return Config from...
func NewConfig(flagCmd CommandLineI) (*Config, error) {
	db := dataBaseConnection{
		Host:     "localhost",
		Port:     7777,
		Table:    "metida",
		User:     "postgres",
		Password: "postgres",
	}

	return &Config{db, flagCmd}, nil
}

// String provides human-readable representation Config.
func (o Config) String() string {
	return fmt.Sprintf(
		`
		Config api:
			data base connect:
			%v

		`,
		o.db,
	)
}

func (o *Config) GetConnectDB() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		o.db.User,
		o.db.Password,
		o.db.Host,
		o.db.Port,
		o.db.Table,
	)
}
