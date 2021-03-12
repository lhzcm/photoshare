package config

type RedisConnConfig struct {
	Address  string `yaml:"address"`
	Dbnum    int    `yaml:"dbnum"`
	Password string `yaml:"password"`
}

type RedisConfig struct {
	User RedisConnConfig `yaml:"user"`
	Msg  RedisConnConfig `yaml:"msg"`
}
