package config

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kisielk/sqlstruct"
)

type DBconf struct {
	Host     string `json:"host" yaml:"host"`
	Database string `json:"database" yaml:"database"`
	Schema   string `json:"schema" yaml:"schema"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Driver   string `json:"driver" yaml:"driver"`
	Port     int    `json:"port" yaml:"port"`
}

func init() {
	sqlstruct.NameMapper = sqlstruct.ToSnakeCase
}

// DBConnection creates a database connection from the config
func (conf *DBconf) DBConnection() (*sql.DB, error) {
	switch {
	case conf.Driver == "mysql":
		return sql.Open("mysql", mysqlConnectionString(conf))
	case conf.Driver == "postgres":
		return sql.Open("postgres", postgresConnectionString(conf))
	}
	return nil, errors.New("Not supported database driver.")
}

// user:password@tcp(localhost:5555)/dbname
// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func mysqlConnectionString(conf *DBconf) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database)
}

// postgresConnectionString
// example: "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
// * dbname - The name of the database to connect to
// * user - The user to sign in as
// * password - The user's password
// * host - The host to connect to. Values that start with / are for unix domain sockets. (default is localhost)
// * port - The port to bind to. (default is 5432)
// * sslmode - Whether or not to use SSL (default is require, this is not the default for libpq)
// * fallback_application_name - An application_name to fall back to if one isn't provided.
// * connect_timeout - Maximum wait for connection, in seconds. Zero or not specified means wait indefinitely.
// * sslcert - Cert file location. The file must contain PEM encoded data.
// * sslkey - Key file location. The file must contain PEM encoded data.
// * sslrootcert - The location of the root certificate file. The file must contain PEM encoded data.
func postgresConnectionString(conf *DBconf) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database)
}
