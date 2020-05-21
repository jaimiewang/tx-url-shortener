package config

import (
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
	BaseURLLength  int    `yaml:"base_url_length"`
	ShortURLPrefix string `yaml:"short_url_prefix"`
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

	return nil
}
