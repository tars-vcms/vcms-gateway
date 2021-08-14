package routes

import (
	"github.com/tars-vcms/vcms-gateway/entity/route"
	"net/http"
	"sync"
)

type HttpRouteManager interface {
	ParseRouteHttp(request *http.Request) *route.HttpRoute

	SubscribeRoute()
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
