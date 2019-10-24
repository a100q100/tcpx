package config

type PostgresNode struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
	Password string `yaml:"password"`
}
type PostgresConfig struct {
	Enable      bool           `yaml:"enable"`
	WriteDBNode PostgresNode   `yaml:"write"`
	ReadDBNodes []PostgresNode `yaml:"read"`
}
