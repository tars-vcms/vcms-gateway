package routes

import (
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/util/rogger"
	"github.com/go-redis/redis/v8"
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
	routesUpdater HttpRouteUpdater
}

func newHttpRouteManagerImpl() *HttpRouteManagerImpl {
	rcfg := rcfgs.GetInstance()
	gateway := &config.GatewayConfig{
		HeaderMap: make(map[string]string),
	}
	if err := rcfg.GetConfig(config.GATEWAY_FILE_NAME, rcfgs.STRUCT, gateway); err != nil {
		panic(err)
	}
	manager := &HttpRouteManagerImpl{
		logger:      tars.GetLogger("CLOG"),
		rcfg:        rcfgs.GetInstance(),
		gatewayName: gateway.Name,
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
	manager.routesUpdater = newHttpRouteUpdater(manager.gatewayName, manager)
	return manager
}

func (h *HttpRouteManagerImpl) SubscribeRoute() {
	h.routesUpdater.Subscribe()
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

func (h *HttpRouteManagerImpl) HandleRoutesUpdate(routes []*route.HttpRoute) {
	h.routesRwMutex.Lock()
	defer h.routesRwMutex.Unlock()
	h.routes = routes
}
