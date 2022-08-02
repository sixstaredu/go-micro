package redis

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type Config struct {
	// 链接地址
	Address string  `mapstructure:"address`
	// 最大活跃链接数
	MaxNum int		`mapstructure:"max_num`
	// 最大闲置连接数
	MaxIdle int		`mapstructure:"max_idle`
	// 闲置连接超时时间
	IdleTimeout int `mapstructure:"idle_timeout`
}


func InitRedisPool(cfg *Config) *redis.Pool {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.Address)
		},
		TestOnBorrow:    nil,
		MaxIdle:         cfg.MaxIdle,
		MaxActive:       cfg.MaxNum,
		IdleTimeout:     time.Second * time.Duration( cfg.IdleTimeout),
		Wait:            false,
		MaxConnLifetime: 0,
	}
	return pool
}