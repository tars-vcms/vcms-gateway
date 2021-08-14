package config

const (
	REDIS_FILE_NAME = "cache.conf"
)

type RedisConfig struct {
	Addr     string `tars:"/cache/Addr"`
	Password string `tars:"/cache/Password"`
	DB       int    `tars:"/cache/DB"`
}
