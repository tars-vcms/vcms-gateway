package rcfg

import (
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/util/conf"
)

type ParseFunc func(content string, dest interface{}) error

type RemoteCfgImpl struct {
	tarsRConf *tars.RConf
}

func (r RemoteCfgImpl) ParseConfig(content string) (*conf.Conf, error) {
	c := conf.New()
	if err := c.InitFromString(content); err != nil {
		return nil, err
	}
	return c, nil
}

func newRemoteCfgImpl() *RemoteCfgImpl {
	cfg := tars.GetServerConfig()
	impl := &RemoteCfgImpl{
		tarsRConf: tars.NewRConf(cfg.App, cfg.Server, cfg.BasePath),
	}
	return impl
}

func (r RemoteCfgImpl) GetConfig(filename string, cfgType CfgType, dest interface{}) error {
	content, err := r.tarsRConf.GetConfig(filename)
	if err != nil {
		return err
	}

	err = getParseFunc(cfgType)(content, dest)
	if err != nil {
		return err
	}
	return nil
}
