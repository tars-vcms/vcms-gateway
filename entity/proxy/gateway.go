package proxy

import (
	"github.com/tars-vcms/vcms-gateway/entity/route"
	"net/http"
)

type HandleResponseHeaderFunc func(resp http.Header) error

type GatewayProxyCallback struct {
	HandleResponseHeader HandleResponseHeaderFunc
}

type GatewayProxy interface {
	DoRequest(routeHttp *route.HttpRoute, w http.ResponseWriter, r *http.Request) (httpCode int, err error)

	SetCallback(cb *GatewayProxyCallback)
}
