package proxymanager

import (
	"github.com/tars-vcms/vcms-gateway/entity/proxy"
	"github.com/tars-vcms/vcms-gateway/entity/route"
	"net/http"
)

type GatewayProxyManager interface {
	GetProxy(httpRoute *route.HttpRoute, r *http.Request) (proxy.GatewayProxy, error)
}

func NewGatewayProxyManager() GatewayProxyManager {
	return newGatewayProxyManagerImpl()
}
