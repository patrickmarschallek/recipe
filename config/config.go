package config

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"path/filepath"
	"reflect"

	"os"

	"strconv"

	yaml "gopkg.in/yaml.v2"
)

// Settings for the service
var Settings Config

// EnvSettings for the service
var envSettings Config

// EnvSettings for the service
var flagSettings Config

//Config service configuration
type Config struct {
	ServiceName string `json:"name" yaml:"name"`
	Host        string `json:"host" yaml:"host"`
	Port        string `json:"port" yaml:"port"`
	BasePath    string `json:"basePath" yaml:"basePath"`
	Persistence *DBconf
}

func init() {
	initFlag()
	initEnv()
}

func initFlag() {
	flagSettings = Config{
		Persistence: &DBconf{},
	}

	// service config flags
	flag.StringVar(&flagSettings.ServiceName, "name", "", "name of the service")
	flag.StringVar(&flagSettings.Host, "host", "", "host address of the service")
	flag.StringVar(&flagSettings.Port, "port", "80", "bind to defined ports")

	//database config flags
	flag.StringVar(&flagSettings.Persistence.Host, "db-host", "", "host address for the database")
	flag.StringVar(&flagSettings.Persistence.Driver, "db-driver", "", "database driver: [mysql,postgres]")
	flag.StringVar(&flagSettings.Persistence.Database, "db-name", "", "database name which will be used")
	flag.StringVar(&flagSettings.Persistence.Password, "db-password", "", "database password for authentification")
	flag.StringVar(&flagSettings.Persistence.User, "db-user", "", "database user for the authentication")
	flag.StringVar(&flagSettings.Persistence.Schema, "db-schema", "", "database schema just for postgres or oracle")
	flag.IntVar(&flagSettings.Persistence.Port, "db-port", 0, "database port which will used for the connection")

	flag.Parse()
}

func initEnv() {
	envSettings = Config{
		Persistence: &DBconf{},
	}

	// service config flags
	envSettings.ServiceName = os.Getenv("RECIPE_NAME")
	envSettings.Host = os.Getenv("RECIPE_HOST")
	envSettings.Port = os.Getenv("RECIPE_PORT")

	//database config flags
	envSettings.Persistence.Host = os.Getenv("RECIPE_DB_HOST")
	envSettings.Persistence.Driver = os.Getenv("RECIPE_DB_DRIVER")
	envSettings.Persistence.Database = os.Getenv("RECIPE_DB_NAME")
	envSettings.Persistence.Password = os.Getenv("RECIPE_DB_PASSWORD")
	envSettings.Persistence.User = os.Getenv("RECIPE_DB_USER")
	envSettings.Persistence.Port, _ = strconv.Atoi(os.Getenv("RECIPE_DB_PORT"))
	envSettings.Persistence.Schema = os.Getenv("RECIPE_DB_SCHEMA")
}

// InitConfig initialize configuration
func InitConfig(path string) (*Config, error) {
	config, err := ReadConfig(path)
	mergeConfigs(&flagSettings, envSettings)
	mergeConfigs(config, flagSettings)
	return config, err
}

// ReadConfig reads the config from the given path
func ReadConfig(path string) (*Config, error) {
	ext := filepath.Ext(path)[1:]
	switch {
	case ext == "yaml" || ext == "yml":
		return ReadYAMLConfig(path)
	case ext == "json":
		return ReadJSONConfig(path)
	}
	return nil, errors.New("unkown config file format")
}

// ReadJSONConfig reads the config from a file
func ReadJSONConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &Settings)
	return &Settings, err
}

// ReadYAMLConfig reads the config from a file
func ReadYAMLConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, &Settings)
	return &Settings, err
}

func mergeConfigs(changed *Config, override Config) {
	mergeStructs(reflect.ValueOf(changed), reflect.ValueOf(override))
}

func mergeStructs(c reflect.Value, o reflect.Value) {
	for i := 0; i < o.NumField(); i++ {
		fieldO := o.Field(i)
		fieldC := c.Elem().Field(i)
		switch fieldO.Interface().(type) {
		case *DBconf:
			mergeStructs(fieldC, fieldO.Elem())
		default:
			if fieldO.Interface() != reflect.Zero(reflect.TypeOf(fieldO.Interface())).Interface() {
				fieldC.Set(reflect.Value(fieldO))
			}
		}
	}
}

func IsEmpty(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
