package config

type ServerConfig struct {
	Port     int    `yaml:"port"`
	Corshost string `yaml:"corshost"`
}
