package internal

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
    Threads  int    `yaml:"Threads"`
    Timeout  int    `yaml:"Timeout"`
    Delay    string `yaml:"Delay"`
    Status   string `yaml:"Status"`
    Rate     int    `yaml:"Rate"`
    Method   string `yaml:"Method"`
    Wordlist string `yaml:"Wordlist"`
    Userlist string `yaml:"Userlist"`
    Passlist string `yaml:"Passlist"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
