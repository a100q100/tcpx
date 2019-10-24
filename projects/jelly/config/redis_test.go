package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestRedisYAML(t *testing.T) {
	var redisCloud = RedisConfig{
		WriteNodeUrl: "redis://localhost:6379",
		ReadNodeUrl: []string{
			"redis://localhost:6379",
			"redis://localhost:6379",
		},
	}

	buf, e := yaml.Marshal(redisCloud)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Println(string(buf))
}
