package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/tars-vcms/vcms-gateway/entity/cache"
	"github.com/tars-vcms/vcms-gateway/entity/route"
	"sync"
	"time"
)

func (h *HttpRouteManagerImpl) routeListenGuard() {
	wg := &sync.WaitGroup{}
	for true {
		wg.Add(1)
		h.logger.Info("[Routes] Route Listen Start")
		go h.routeListen(wg)
		wg.Wait()
	}
}

func (h *HttpRouteManagerImpl) routeListen(wg *sync.WaitGroup) {

	ctx := context.Background()
	// 订阅频道为 网关名 + 固定后缀组成的频道
	ps := h.redis.Subscribe(ctx, h.gatewayName+cache.MQ_SUFFIX)
	defer func(ps *redis.PubSub) {
		err := ps.Close()
		if err != nil {
			h.logger.Error("[Routes] Redis PubSub Close failed %v", err.Error())
		}
		wg.Done()
	}(ps)
	_, err := ps.Receive(ctx)
	if err != nil {
		h.logger.Error("[Routes] Redis PubSub Receive failed %s", err.Error())
		return
	}
	for msg := range ps.Channel() {
		mQRouteCmd := &cache.MQRouteCmd{}
		err := json.Unmarshal([]byte(msg.Payload), mQRouteCmd)
		if err != nil {
			h.logger.Error("[Routes] Unmarshal MQRouteCmd failed %s", err.Error())
			continue
		}
		switch mQRouteCmd.CMD {
		case cache.CMD_ROUTE_UPDATE:
			if err := h.loadRoute(mQRouteCmd.Payload); err != nil {
				h.logger.Error("[Routes] Load Route failed %v", err.Error())
			}
			break
		default:
			h.logger.Info("[Routes] Undefined MQRouteCmd %s", mQRouteCmd.CMD)

		}
		fmt.Printf("channel=%s message=%s\n", msg.Channel, msg.Payload)
	}
}

func (h *HttpRouteManagerImpl) loadRoute(version string) error {
	ctx := context.Background()
	stringCmd := h.redis.Get(ctx, h.gatewayName+cache.ROUTE_VERSION_SUFFIX+"_"+version)
	if stringCmd.Err() != nil {
		return stringCmd.Err()
	}
	h.routesRwMutex.Lock()
	defer h.routesRwMutex.Unlock()
	var r []*route.HttpRoute
	err := json.Unmarshal([]byte(stringCmd.Val()), r)
	if err != nil {
		return err
	}
	h.routes = r
	return nil
}

func (h *HttpRouteManagerImpl) pubRequireRoute() error {
	ctx := context.Background()
	mQRouteCmd := &cache.MQRouteCmd{
		CMD:       cache.CMD_ROUTE_REQUIRE,
		Timestamp: time.Now().Unix(),
	}
	b, err := json.Marshal(mQRouteCmd)
	if err != nil {
		return err
	}
	return h.redis.Publish(ctx, h.gatewayName+cache.MQ_SUFFIX, b).Err()
}
