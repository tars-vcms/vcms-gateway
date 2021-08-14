package rcfgs

import (
	"github.com/go-redis/redis/v8"
	"sync"
)

type CfgType uint8

const (
	JSON   CfgType = 1
	YAML   CfgType = 2
	TEXT   CfgType = 3
	STRUCT CfgType = 4
)

type RemoteCfg interface {
	GetConfig(filename string, cfgType CfgType, dest interface{}) error

	GetRedisClient() *redis.Client
}

var remoteCfg RemoteCfg
var once sync.Once

func GetInstance() RemoteCfg {
	once.Do(func() {
		remoteCfg = newRemoteCfgImpl()
	})
	return remoteCfg
}
