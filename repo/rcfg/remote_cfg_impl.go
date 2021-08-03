package rcfg

import (
	"encoding/json"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/util/conf"
)

type ParseFunc func(content string, dest interface{}) error

type RemoteCfgImpl struct {
	tarsRConf *tars.RConf
	parseMap  map[CfgType]ParseFunc
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
	impl.parseMap = map[CfgType]ParseFunc{
		JSON: impl.parseJSON,
		YAML: impl.parseYAML,
		TEXT: impl.parseText,
	}
	return impl
}

func (r RemoteCfgImpl) GetConfig(filename string, cfgType CfgType, dest interface{}) error {
	content, err := r.tarsRConf.GetConfig(filename)
	if err != nil {
		return err
	}

	err = r.getParseFunc(cfgType)(content, dest)
	if err != nil {
		return err
	}
	return nil
}

func (r RemoteCfgImpl) getParseFunc(cfgType CfgType) ParseFunc {
	return r.parseMap[cfgType]
}

func (r RemoteCfgImpl) parseJSON(content string, dest interface{}) error {
	return json.Unmarshal([]byte(content), dest)
}

func (r RemoteCfgImpl) parseYAML(content string, dest interface{}) error {
	return json.Unmarshal([]byte(content), dest)
}

func (r RemoteCfgImpl) parseText(content string, dest interface{}) error {
	*(dest).(*string) = content
	return nil
}
