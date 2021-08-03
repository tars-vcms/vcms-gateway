package routemanager

import (
	"github.com/tars-vcms/vcms-gateway/entity/route"
	"net/http"
)

type HttpRouteManager interface {
	ParseRouteHttp(request *http.Request) *route.HttpRoute
}

func NewHttpRouteManager() HttpRouteManager {
	return newHttpRouteManagerImpl()
}
