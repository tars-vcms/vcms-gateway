package main

import (
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/tars-vcms/vcms-common/errs"
	"github.com/tars-vcms/vcms-gateway/entity/proxy"
	"github.com/tars-vcms/vcms-gateway/repo/auth"
	"github.com/tars-vcms/vcms-gateway/repo/header"
	"github.com/tars-vcms/vcms-gateway/repo/proxies"
	"github.com/tars-vcms/vcms-gateway/repo/routes"
	"net/http"
	"strconv"
)

func newProxyHttpImp() *ProxyHttpImp {
	mux := &tars.TarsHttpMux{}
	proxyHttp := &ProxyHttpImp{
		mux:    mux,
		route:  routes.GetInstance(),
		auth:   auth.NewHttpAuth(),
		proxy:  proxies.NewGatewayProxyManager(),
		header: header.NewHttpHeader(),
	}
	mux.HandleFunc("/", proxyHttp.requestHandler)
	return proxyHttp
}

type ProxyHttpImp struct {
	mux    *tars.TarsHttpMux
	route  routes.HttpRouteManager
	auth   auth.HttpAuth
	proxy  proxies.GatewayProxyManager
	header header.HttpHeader
}

func (p ProxyHttpImp) handleError(w http.ResponseWriter, httpCode int, err error) {
	h := w.Header()
	p.header.InjectCommonHeader(h)
	code := errs.Code(err)
	msg := errs.Msg(err)
	h.Set("tars-ret", strconv.Itoa(code))
	h.Set("tars-msg", msg)
	w.WriteHeader(httpCode)
}

func (p ProxyHttpImp) GetTarsHttpMux() *tars.TarsHttpMux {
	return p.mux
}

func (p ProxyHttpImp) requestHandler(w http.ResponseWriter, r *http.Request) {

	if httpCode, err := p.header.CheckRequestHeader(r); err != nil {
		p.handleError(w, httpCode, err)
		return
	}

	routeHttp := p.route.ParseRouteHttp(r)
	// 判断是否找到路由，没找到就是404
	if routeHttp == nil {
		p.handleError(w, 404, nil)
		return
	}
	// 判断是否有权限访问路由，没有权限返回401
	if err := p.auth.ChallengeHttpAuth(r, w, routeHttp); err != nil {
		p.handleError(w, 401, err)
		return
	}

	gProxy, err := p.proxy.GetProxy(routeHttp, r)
	if err != nil {
		p.handleError(w, 502, err)
		return
	}

	gProxy.SetCallback(&proxy.GatewayProxyCallback{
		HandleResponseHeader: func(header http.Header) error {
			p.header.InjectCommonHeader(header)
			return nil
		},
	})

	if httpCode, err := gProxy.DoRequest(routeHttp, w, r); err != nil {
		p.handleError(w, httpCode, err)
		return
	}
}
