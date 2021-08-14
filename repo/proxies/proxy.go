package proxies

import (
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/basef"
	"github.com/tars-vcms/vcms-common/errs"
	"github.com/tars-vcms/vcms-gateway/entity/errcode"
	"github.com/tars-vcms/vcms-gateway/entity/proxy"
	"github.com/tars-vcms/vcms-gateway/entity/route"
	"net/http"
	"net/url"
)

type GatewayProxyManagerImpl struct {
	comm *tars.Communicator
}

func (t GatewayProxyManagerImpl) getTupProxy(httpRoute *route.HttpRoute, r *http.Request) (proxy.GatewayProxy, error) {
	if len(httpRoute.ServantName) == 0 {
		return nil, errs.New(errcode.RetServantNameErr, "ServantName Error")
	}
	sProxy := tars.NewServantProxy(t.comm, httpRoute.ServantName)
	sProxy.TarsSetVersion(basef.JSONVERSION)
	return proxy.NewTupProxy(sProxy), nil
}

func (t GatewayProxyManagerImpl) getReserveProxy(httpRoute *route.HttpRoute, r *http.Request) (proxy.GatewayProxy, error) {
	remote, err := url.Parse(httpRoute.ServantName)
	if err != nil {
		return nil, err
	}

	return proxy.NewReserveProxy(remote), nil
}

func (t GatewayProxyManagerImpl) GetProxy(httpRoute *route.HttpRoute, r *http.Request) (proxy.GatewayProxy, error) {
	var p proxy.GatewayProxy
	var err error
	switch httpRoute.Type {
	case route.TARS_SERVANT:
		p, err = t.getTupProxy(httpRoute, r)
		break
	case route.RESERVE_PROXY_SERVANT:
		p, err = t.getReserveProxy(httpRoute, r)
		break
	}
	return p, err
}

func newGatewayProxyManagerImpl() *GatewayProxyManagerImpl {
	return &GatewayProxyManagerImpl{
		comm: tars.NewCommunicator(),
	}
}
