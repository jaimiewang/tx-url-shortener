package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	Database struct {
		Engine   string `yaml:"engine"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Hostname string `yaml:"hostname"`
		Port     uint16 `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	ListenAddress  string `yaml:"listen_address"`
	Secret         string `yaml:"secret"`
	BaseCodeLength int    `yaml:"base_code_length"`
	ShortURLPrefix string `yaml:"short_url_prefix"`
	CacheSize      int    `yaml:"cache_size"`
}

var Config = &config{}

func LoadConfig(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, Config)
	if err != nil {
		return err
	}

	if Config.CacheSize < 1 {
		return errors.New("cache_size must be positive")
	}

	if Config.BaseCodeLength > 11 {
		return errors.New("base_code_length must be smaller than or equal to 11")
	}

	return nil
}
