package redis

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var Pool *redis.Pool

func Init(dataSource string) {
	Pool = &redis.Pool{
		MaxIdle: 8,
		MaxActive: 6,
		IdleTimeout: 200*time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", dataSource)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}

func Int(reply interface{}, err error) (int, error) {
	return redis.Int(reply, err)
}
