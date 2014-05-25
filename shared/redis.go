package shared

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var RedisPool *redis.Pool

func InitializeRedisPool() {
	if RedisPool == nil {
		redisHost := Config.Redis.Server
		if len(redisHost) == 0 {
			redisHost = ":6379"
		}

		RedisPool = newPool(redisHost, Config.Redis.Password)
		if RedisPool == nil {
			log.Fatal("Unable to connect to Redis!")
		} else {
			log.Println("Connected to Redis")
		}
	}
}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 2 * time.Second,
		MaxActive:   50,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if len(password) != 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
