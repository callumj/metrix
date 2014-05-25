package shared

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
)

var Config Configuration

type Configuration struct {
	Sentry string
	Redis  RedisConfig
	Listen string
}

type RedisConfig struct {
	Server   string
	Password string
}

func LoadConfig(path string) bool {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Unable to find configuration %v (%v)\r\n", path, err)
		return false
	}

	err = yaml.Unmarshal([]byte(fileContents), &Config)
	if err != nil {
		log.Printf("Unable to parse configuration %v\r\n", path)
		return false
	}

	return true
}

func ExplainConfig() {
	log.Printf("Sentry DSN: %s", Config.Sentry)
	log.Printf("Redis host: %s", Config.Redis.Server)
	log.Printf("Redis password: %s", Config.Redis.Password)
}
