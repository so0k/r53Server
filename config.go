package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config specifies the r53server config
type Config struct {
	Roles []RoleConfig `yaml:"roles,omitempty"`
}

// RoleConfig specifies the configuration of a role
type RoleConfig struct {
	RoleArn string   `yaml:"roleArn,omitempty"`
	Zones   []string `yaml:"zones,omitempty"`
}

func readConfig(yamlfile string) (*Config, error) {
	c := Config{}
	data, err := ioutil.ReadFile(yamlfile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
