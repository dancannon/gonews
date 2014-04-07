package infrastructure

import (
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/dancannon/gonews/core/config"
)

var (
	redis_pool *redis.Pool
)

func InitRedis(conf config.Redis) {
	redis_pool = &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: conf.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(conf.Network, conf.Address)
			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func Redis() *redis.Pool {
	return redis_pool
}
