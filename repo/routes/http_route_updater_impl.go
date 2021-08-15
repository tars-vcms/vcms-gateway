package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/util/rogger"
	"github.com/go-redis/redis/v8"
	"github.com/tars-vcms/vcms-gateway/entity/cache"
	"github.com/tars-vcms/vcms-gateway/entity/route"
	"github.com/tars-vcms/vcms-gateway/repo/rcfgs"
	"sync"
	"time"
)

func newHttpRouteUpdaterImpl(name string, routeManager HttpRouteManager) HttpRouteUpdater {
	return &HttpRouteUpdaterImpl{
		callback: routeManager.HandleRoutesUpdate,
		logger:   tars.GetLogger("CLOG"),
		name:     name,
		redis:    rcfgs.GetInstance().GetRedisClient(),
	}
}

type HttpRouteUpdaterImpl struct {
	callback func(routes []*route.HttpRoute)
	wg       sync.WaitGroup
	logger   *rogger.Logger
	name     string
	redis    *redis.Client
}

func (h *HttpRouteUpdaterImpl) Subscribe() {
	version, err := h.getRouteVersion()
	ok := err == nil
	if ok {
		if err := h.updateRoutes(version); err != nil {
			h.logger.Error(fmt.Sprintf("[Routes] Load Route failed %v", err.Error()))
			// 从缓存加载失败，尝试向网关重新获取
			ok = false
		}
	} else if err != nil {
		h.logger.Error(fmt.Sprintf("[Routes] Get Route Version failed %v", err.Error()))
	}
	h.Listen()
	if !ok {
		h.logger.Info("[Routes] Try to Pub Require Route Cmd ")
		if err := h.PubRequireRoute(); err != nil {
			h.logger.Error(fmt.Sprintf("[Routes] PubRequireRoute failed %v", err.Error()))
		}
	}
}

func (h *HttpRouteUpdaterImpl) Listen() {
	go func() {
		for true {
			h.wg.Add(1)
			h.logger.Info("[Routes] Route Listen Start")
			go h.listenHandler()
			h.wg.Wait()
		}
	}()
}

func (h *HttpRouteUpdaterImpl) listenHandler() {
	ctx := context.Background()
	// 订阅频道为 网关名 + 固定后缀组成的频道
	ps := h.redis.Subscribe(ctx, h.name+cache.MQ_SUFFIX)
	defer func(ps *redis.PubSub) {
		err := ps.Close()
		if err != nil {
			h.logger.Error(fmt.Sprintf("[Routes] Redis PubSub Close failed %v", err.Error()))
		}
		h.wg.Done()
	}(ps)
	_, err := ps.Receive(ctx)
	if err != nil {
		h.logger.Error(fmt.Sprintf("[Routes] Redis PubSub Receive failed %s", err.Error()))
		return
	}
	for msg := range ps.Channel() {
		mQRouteCmd := &cache.MQRouteCmd{}
		err := json.Unmarshal([]byte(msg.Payload), mQRouteCmd)
		if err != nil {
			h.logger.Error(fmt.Sprintf("[Routes] Unmarshal MQRouteCmd failed %s", err.Error()))
			continue
		}
		switch mQRouteCmd.CMD {
		case cache.CMD_ROUTE_UPDATE:
			if err := h.updateRoutes(mQRouteCmd.Payload); err != nil {
				h.logger.Error(fmt.Sprintf("[Routes] Load Route failed %v", err.Error()))
			}
			break
		default:
			h.logger.Info(fmt.Sprintf("[Routes] Undefined MQRouteCmd %v", mQRouteCmd.CMD))

		}
		//fmt.Printf("channel=%s message=%s\n", msg.Channel, msg.Payload)
	}
}

func (h *HttpRouteUpdaterImpl) updateRoutes(version string) error {
	ctx := context.Background()
	stringCmd := h.redis.Get(ctx, h.name+cache.ROUTE_VERSION_SUFFIX+"_"+version)
	if stringCmd.Err() != nil {
		return stringCmd.Err()
	}
	r := make([]*route.HttpRoute, 0)
	err := json.Unmarshal([]byte(stringCmd.Val()), &r)
	if err != nil {
		return err
	}
	h.callback(r)
	return nil
}

func (h *HttpRouteUpdaterImpl) getRouteVersion() (string, error) {
	ctx := context.Background()
	cmd := h.redis.Get(ctx, h.name+cache.ROUTE_VERSION_SUFFIX)
	return cmd.Val(), cmd.Err()
}

func (h *HttpRouteUpdaterImpl) PubRequireRoute() error {
	ctx := context.Background()
	mQRouteCmd := &cache.MQRouteCmd{
		CMD:       cache.CMD_ROUTE_REQUIRE,
		Timestamp: time.Now().Unix(),
	}
	b, err := json.Marshal(mQRouteCmd)
	if err != nil {
		return err
	}
	return h.redis.Publish(ctx, h.name+cache.MQ_SUFFIX, b).Err()
}
