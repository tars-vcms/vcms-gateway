package routes

import (
	"github.com/tars-vcms/vcms-gateway/entity/route"
	"net/http"
	"sync"
)

type HttpRouteManager interface {
	ParseRouteHttp(request *http.Request) *route.HttpRoute

	SubscribeRoute()

	HandleRoutesUpdate(routes []*route.HttpRoute)
}

type HttpRouteUpdater interface {
	Subscribe()

	Listen()
}

var once sync.Once
var manager HttpRouteManager

func GetInstance() HttpRouteManager {
	once.Do(func() {
		manager = newHttpRouteManagerImpl()
		manager.SubscribeRoute()
	})
	return manager
}

func newHttpRouteUpdater(name string, routeManager HttpRouteManager) HttpRouteUpdater {
	return newHttpRouteUpdaterImpl(name, routeManager)
}
