package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestPostgresDB(t *testing.T) {
	b, e := json.MarshalIndent(PGDB, "  ", "  ")
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Println(string(b))
}

func TestRedisDB(t *testing.T) {
	b, e := json.MarshalIndent(RedisDB, "  ", "  ")
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Println(string(b))
}

func TestPostgresYAML(t *testing.T) {
	var write = PostgresNode{
		Host:     "localhost",
		User:     "postgres",
		Dbname:   "jelly",
		Password: "123",
		Sslmode:  "disable",
	}
	var read1 = PostgresNode{
		Host:     "localhost",
		User:     "postgres",
		Dbname:   "jelly",
		Password: "123",
		Sslmode:  "disable",
	}
	var read2 = PostgresNode{
		Host:     "localhost",
		User:     "postgres",
		Dbname:   "jelly",
		Password: "123",
		Sslmode:  "disable",
	}

	pgConfig := PostgresConfig{
		WriteDBNode: write,
		ReadDBNodes: []PostgresNode{
			read1,
			read2,
		},
	}
	buf, e := yaml.Marshal(pgConfig)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Println(string(buf))
}
