package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Host           string
	Port           string
	MaxConnections int
}

func Connect(cfg Config) *redis.Pool {
	redisAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	return redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", redisAddr)
	}, cfg.MaxConnections)
}
