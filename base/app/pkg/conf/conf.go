package conf

import (
	"time"
)

var (
	defaultConf *Config
)

func Default() *Config {
	return defaultConf
}

type Config struct {
	DB    *DB
	Redis *Redis
	Dfs   *Dfs
	Etcd  *Etcd
}

type DB struct {
	Args         string
	IdleTimeout  time.Duration
	QueryTimeout time.Duration
	ExecTimeout  time.Duration
	TranTimeout  time.Duration
}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type Dfs struct {
	Host  string
	Group string
}

type Etcd struct {
	Addr string
}

func init() {
	defaultConf = &Config{
		DB: &DB{
			Args: "tiktok:123456@tcp(127.0.0.1:3306)/db_tiktok?charset=utf8&parseTime=true&loc=Asia%2FChongqing",
		},
		Redis: &Redis{
			Host:        "127.0.0.1:6379",
			Password:    "",
			MaxIdle:     500,
			MaxActive:   200,
			IdleTimeout: 500,
		},
		Dfs: &Dfs{
			Host:  "192.168.43.6:8086",
			Group: "group1",
		},
		Etcd: &Etcd{
			Addr: "172.28.5.1:2379",
		},
	}
}
