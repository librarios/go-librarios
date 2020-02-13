package app

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	DefaultPort = 8080
)

type Config struct {
	Port    int                               `yaml:"port"`
	Plugins map[string]map[string]interface{} `yaml:"plugins"`
	DB      map[string]interface{}            `yaml:"db"`
}

// LoadConfigFile reads configuration file
func LoadConfigFile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return LoadConfigBytes(data)
}

// LoadConfigBytes reads configuration from byte array
func LoadConfigBytes(data []byte) (*Config, error) {
	config := Config{
		Port: DefaultPort,
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
