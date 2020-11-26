package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

const configfile = "./config.yaml"

var Configs *Config

func init() {
	yamlfile, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Fatalln(err)
	}
	Configs = &Config{}
	if err = yaml.Unmarshal(yamlfile, Configs); err != nil {
		log.Fatalln(err)
	}
}

type Config struct {
	Token TokenConfig `yaml:"token"`
	Mssql MssqlConfig `yaml:"mssql"`
	Redis RedisConfig `yaml:"redis"`
}
