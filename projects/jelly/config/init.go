package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
)

var fileDir string
var DBMode string

// DB
var (
	PGDB    PostgresConfig
	RedisDB RedisConfig
)

func init() {
	flag.StringVar(&fileDir, "fd", "G:\\go_workspace\\GOPATH\\src\\tcpx\\projects\\jelly\\config\\files", "-fd <cfg/file/dir>")
	flag.Parse()

    initPostgresDB()
	initRedisDB()
}

// init postgres
func initPostgresDB() {
	filePath := path.Join(fileDir, "postgres.yaml")
	buf, e := ioutil.ReadFile(filePath)
	if e != nil {
		panic(e)
	}

	if e := yaml.Unmarshal(buf, &PGDB); e != nil {
		panic(e)
	}
}

// init redis
func initRedisDB() {
	filePath := path.Join(fileDir, "redis.yaml")
	buf, e := ioutil.ReadFile(filePath)
	if e != nil {
		panic(e)
	}

	if e := yaml.Unmarshal(buf, &RedisDB); e != nil {
		panic(e)
	}
}
