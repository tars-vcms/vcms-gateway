package rcfgs

import (
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/util/rogger"
	"github.com/go-redis/redis/v8"
	"github.com/tars-vcms/vcms-gateway/entity/config"
	"sync"
)

type ParseFunc func(content string, dest interface{}) error

type RemoteCfgImpl struct {
	tarsRConf *tars.RConf
	cacheMap  map[string]string
	logger    *rogger.Logger
	redis     *redis.Client
	redisOnce sync.Once
}

func (r *RemoteCfgImpl) GetRedisClient() *redis.Client {
	r.redisOnce.Do(func() {
		options, err := r.getRedisOptions()
		if err != nil {
			r.logger.Error(fmt.Sprintf("[RemoteCfg] GetRedis Options failed:%v\n", err.Error()))
			panic(err)
		}
		r.redis = redis.NewClient(options)
	})
	return r.redis
}

func newRemoteCfgImpl() *RemoteCfgImpl {
	cfg := tars.GetServerConfig()
	impl := &RemoteCfgImpl{
		tarsRConf: tars.NewRConf(cfg.App, cfg.Server, cfg.BasePath),
		cacheMap:  make(map[string]string),
		logger:    tars.GetLogger("CLOG"),
	}
	return impl
}

func (r *RemoteCfgImpl) GetConfig(filename string, cfgType CfgType, dest interface{}) error {
	var content string
	value, ok := r.cacheMap[filename]
	if ok {
		content = value
	} else {
		value, err := r.tarsRConf.GetConfig(filename)
		if err != nil {
			return err
		}
		content = value
		r.cacheMap[filename] = content
	}

	err := getParseFunc(cfgType)(content, dest)
	if err != nil {
		return err
	}
	return nil
}

func (r *RemoteCfgImpl) getRedisOptions() (*redis.Options, error) {
	redisCfg := &config.RedisConfig{}
	if err := r.GetConfig(config.REDIS_FILE_NAME, STRUCT, redisCfg); err != nil {
		return nil, err
	}
	redisOptions := &redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	}
	return redisOptions, nil
}
