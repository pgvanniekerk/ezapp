package config

import (
	"strconv"
)

// DBConf contains configuration for connecting to a PostgreSQL database.
// This struct demonstrates how to define configuration for external resources
// that your application depends on.
//
// Each field has an `envvar` tag that specifies the environment variable name
// to use for that field, and a `default` tag that specifies the default value
// to use if the environment variable is not set.
//
// In a real application, you would typically load this configuration from
// environment variables using a library like github.com/kelseyhightower/envconfig.
type DBConf struct {
	// Host is the database server hostname or IP address
	// Environment variable: DB_HOST
	// Default: localhost
	Host string `envvar:"DB_HOST" default:"localhost"`

	// Port is the database server port
	// Environment variable: DB_PORT
	// Default: 5432 (standard PostgreSQL port)
	Port int `envvar:"DB_PORT" default:"5432"`

	// User is the database username
	// Environment variable: DB_USER
	// Default: postgres
	User string `envvar:"DB_USER" default:"postgres"`

	// Password is the database password
	// Environment variable: DB_PASSWORD
	// Default: postgres
	Password string `envvar:"DB_PASSWORD" default:"postgres"`

	// DBName is the name of the database to connect to
	// Environment variable: DB_NAME
	// Default: postgres
	DBName string `envvar:"DB_NAME" default:"postgres"`

	// SSLMode is the SSL mode to use for the connection
	// Environment variable: DB_SSL_MODE
	// Default: disable
	// Possible values: disable, require, verify-ca, verify-full
	SSLMode string `envvar:"DB_SSL_MODE" default:"disable"`
}

// GetConnectionString returns a formatted connection string for the database.
// This method formats the configuration fields into a connection string that
// can be used with the database/sql package and the PostgreSQL driver.
//
// Returns:
//   - A connection string in the format required by the PostgreSQL driver
//
// Example:
//
//	dbConf := config.DBConf{
//	    Host:     "localhost",
//	    Port:     5432,
//	    User:     "myuser",
//	    Password: "mypassword",
//	    DBName:   "mydb",
//	    SSLMode:  "disable",
//	}
//	connStr := dbConf.GetConnectionString()
//	db, err := sql.Open("postgres", connStr)
func (c DBConf) GetConnectionString() string {
	return "host=" + c.Host +
		" port=" + strconv.Itoa(c.Port) +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}
