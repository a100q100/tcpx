package config

type RedisConfig struct {
	AsDB         bool     `yaml:"asdb"`
	Enable       bool     `yaml:"enable"`
	WriteNodeUrl string   `yaml:"write"`
	ReadNodeUrl  []string `yaml:"read"`
}
