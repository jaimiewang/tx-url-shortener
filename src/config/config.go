package config

import (
	"fmt"
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
	BaseCodeLength     int    `yaml:"base_code_length"`
	ShortURLPrefix     string `yaml:"short_url_prefix"`
	PrettyAPIResponses bool   `yaml:"pretty_api_responses"`
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

	if Config.BaseCodeLength > 11 {
		return fmt.Errorf("base_code_length must be smaller than or equal to 11, got: %d", Config.BaseCodeLength)
	}

	return nil
}
