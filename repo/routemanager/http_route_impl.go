package routemanager

import (
	"github.com/tars-vcms/vcms-gateway/entity/route"
	"net/http"
	"strings"
)

type HttpRouteManagerImpl struct {
	routes []*route.HttpRoute
}

func newHttpRouteManagerImpl() *HttpRouteManagerImpl {
	return &HttpRouteManagerImpl{
		routes: []*route.HttpRoute{
			{
				Path:        "hello",
				Type:        route.TarsServant,
				ServantName: "StressTest.EchoTestServer.EchoTestObj@tcp -h 127.0.0.1 -p 10017 -t 60000",
				FuncName:    "echo",
				//Children: []*routemanager.HttpRoute{
				//	{
				//		Path: "world",
				//		Type: routemanager.TarsServant,
				//		ServantName: "StressTest.EchoTestServer.EchoTestObj@tcp -h 127.0.0.1 -p 10017 -t 60000",
				//		FuncName: "echo",
				//	},
				//},
			},
			{
				Path:        "world",
				Type:        route.ReverseProxyServant,
				ServantName: "http://localhost:10015/hello",
			},
		},
	}
}

func (h HttpRouteManagerImpl) ParseRouteHttp(request *http.Request) *route.HttpRoute {
	paths := h.splitPath(request.URL.Path)
	return h.parseRouteTree(paths, h.routes)
}

func (h HttpRouteManagerImpl) splitPath(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}

func (h HttpRouteManagerImpl) parseRouteTree(paths []string, routes []*route.HttpRoute) *route.HttpRoute {
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
