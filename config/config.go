package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Settings for the service
var Settings Config

//Config service configuration
type Config struct {
	ServiceName string `json:"name" yaml:"name"`
	Host        string `json:"host" yaml:"host"`
	Port        string `json:"port" yaml:"port"`
	BasePath    string `json:"basePath" yaml:"basePath"`
	Persistence *DBconf
}

func init() {

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
