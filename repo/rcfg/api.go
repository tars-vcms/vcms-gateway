package rcfg

import (
	"github.com/TarsCloud/TarsGo/tars/util/conf"
	"sync"
)

type CfgType uint8

const (
	JSON CfgType = 1
	YAML CfgType = 2
	TEXT CfgType = 3
)

type RemoteCfg interface {
	GetConfig(filename string, cfgType CfgType, dest interface{}) error

	ParseConfig(content string) (*conf.Conf, error)
}

var remoteCfg RemoteCfg
var once sync.Once

func GetInstance() RemoteCfg {
	once.Do(func() {
		remoteCfg = newRemoteCfgImpl()
	})
	return remoteCfg
}
