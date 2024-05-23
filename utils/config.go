package utils

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Env     string `yaml:"env"`
	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
	Database struct {
		ConnectionUrl string `yaml:"connectionUrl"`
	} `yaml:"database"`
}

func LoadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	config := &Config{}
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
