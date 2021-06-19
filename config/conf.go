package config

import (
	"os"
)

const (
	DefaultDBType = "mongodb"
	DefaultDBUser = "admin"
	DefaultDBPass = "example"
	DefaultDBHost = "localhost"
	DefaultDBPort = "27017"
	DefaultDBName = "test"
)

type DBConfig struct {
	DBType string
	DBUser string
	DBPass string
	DBHost string
	DBPort string

	// DBName is the name of the db that you are going to use.
	DBName string
}

// ReadDBConfig parses environment variables. Based on the existing env. variables
// it returns configuration struct to connect a database.
func ReadDBConfig() DBConfig {
	dbc := DBConfig{
		DefaultDBType,
		DefaultDBUser,
		DefaultDBPass,
		DefaultDBHost,
		DefaultDBPort,
		DefaultDBName,
	}

	if v := os.Getenv("DBType"); v != "" {
		dbc.DBType = v
	}

	if v := os.Getenv("DBUser"); v != "" {
		dbc.DBUser = v
	}

	if v := os.Getenv("DBPass"); v != "" {
		dbc.DBPass = v
	}

	if v := os.Getenv("DBHost"); v != "" {
		dbc.DBHost = v
	}

	if v := os.Getenv("DBPort"); v != "" {
		dbc.DBPort = v
	}

	if v := os.Getenv("DBName"); v != "" {
		dbc.DBName = v
	}
	return dbc
}
