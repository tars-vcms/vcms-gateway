package routes

import (
	"context"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/util/rogger"
	"github.com/go-redis/redis/v8"
	"github.com/tars-vcms/vcms-gateway/entity/cache"
	"github.com/tars-vcms/vcms-gateway/entity/config"
	"github.com/tars-vcms/vcms-gateway/entity/route"
	"github.com/tars-vcms/vcms-gateway/repo/rcfgs"
	"net/http"
	"strings"
	"sync"
)

type HttpRouteManagerImpl struct {
	gatewayName   string
	routes        []*route.HttpRoute
	routesRwMutex sync.RWMutex
	redis         *redis.Client
	logger        *rogger.Logger
	rcfg          rcfgs.RemoteCfg
}

func newHttpRouteManagerImpl() *HttpRouteManagerImpl {
	manager := &HttpRouteManagerImpl{
		redis:  rcfgs.GetInstance().GetRedisClient(),
		logger: tars.GetLogger("CLOG"),
		rcfg:   rcfgs.GetInstance(),
		routes: []*route.HttpRoute{
			{
				Path:        "world",
				Type:        route.RESERVE_PROXY_SERVANT,
				ServantName: "http://localhost:10015/hello",
			},
			{
				Path:        "CreateProject",
				Type:        route.TARS_SERVANT,
				ServantName: "vcms.projectmanager.ProjectManagerObj@tcp -h 127.0.0.1 -p 10017 -t 60000",
				FuncName:    "CreateProject",
				InputName:   "input",
				OutputName:  "output",
			},
			{
				Path:        "GetProjects",
				Type:        route.TARS_SERVANT,
				ServantName: "vcms.projectmanager.ProjectManagerObj@tcp -h 127.0.0.1 -p 10017 -t 60000",
				FuncName:    "GetProjects",
				InputName:   "input",
				OutputName:  "output",
			},
		},
	}
	gateway := &config.GatewayConfig{
		HeaderMap: make(map[string]string),
	}
	if err := manager.rcfg.GetConfig(config.GATEWAY_FILE_NAME, rcfgs.STRUCT, gateway); err != nil {
		manager.logger.Error("[Routes] Get GatewayConfig failed %v", err.Error())
		panic(err)
	}
	manager.gatewayName = gateway.Name
	return manager
}

func (h *HttpRouteManagerImpl) SubscribeRoute() {
	ctx := context.Background()
	cmd := h.redis.Get(ctx, h.gatewayName+cache.ROUTE_VERSION_SUFFIX)
	ok := cmd.Err() == nil
	if !ok {
		h.logger.Error("[Routes] Get Route Version failed %v", cmd.Err().Error())
	} else {
		// 首先先预加载一次
		if err := h.loadRoute(cmd.Val()); err != nil {
			h.logger.Error("[Routes] Load Route failed %v", err.Error())
			// 从缓存加载失败，尝试向网关重新获取
			ok = false
		}
	}
	go h.routeListenGuard()
	if !ok {
		h.logger.Info("[Routes] Try to Pub Require Route Cmd ")
		if err := h.pubRequireRoute(); err != nil {
			h.logger.Error("[Routes] PubRequireRoute failed %v", err.Error())
		}
	}
}

func (h *HttpRouteManagerImpl) ParseRouteHttp(request *http.Request) *route.HttpRoute {
	paths := h.splitPath(request.URL.Path)
	return h.parseRouteTree(paths, h.routes)
}

func (h *HttpRouteManagerImpl) splitPath(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}

func (h *HttpRouteManagerImpl) parseRouteTree(paths []string, routes []*route.HttpRoute) *route.HttpRoute {
	h.routesRwMutex.RLock()
	defer h.routesRwMutex.RUnlock()
	if routes == nil {
		return nil
	}
	for _, r1 := range routes {
		path := paths[0]
		if path == r1.Path {
			if len(paths) <= 1 {
				return r1
			}
			if r2 := h.parseRouteTree(paths[1:], r1.Children); r2 != nil {
				return r2
			}
		}

	}

	return nil
}
