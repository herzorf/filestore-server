package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	pool      *redis.Pool
	redisHost = "127.0.0.1:6379"
	redisPass = ""
)

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			//	打开链接
			context, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Println("redis 链接出错", err)
				return nil, err
			}
			if _, err := context.Do("AUTH", redisPass); err != nil {
				_ = context.Close()
				return nil, err
			}
			return context, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
func init() {
	pool = newRedisPool()
}

func RedisPool() *redis.Pool {
	return pool
}
