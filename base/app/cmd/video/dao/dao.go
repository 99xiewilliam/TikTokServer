package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"go-tiktok/app/pkg/conf"
	"time"
)

type Dao struct {
	c     *conf.Config
	db    *gorm.DB
	redis *redis.Pool
}

func New(c *conf.Config) *Dao {
	d := &Dao{
		c:     c,
		db:    open(c),
		redis: newPool(c),
	}
	return d
}

func open(c *conf.Config) *gorm.DB {
	db, err := gorm.Open("mysql", c.DB.Args)
	if err != nil {
		panic(err)
	}
	return db
}

func newPool(c *conf.Config) *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle:     c.Redis.MaxIdle,
		MaxActive:   c.Redis.MaxActive,
		IdleTimeout: c.Redis.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", c.Redis.Host)
			if err != nil {
				return nil, err
			}
			if c.Redis.Password != "" {
				if _, err := conn.Do("AUTH", c.Redis.Password); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return conn, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return redisPool
}
