package main

import (
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/tars-vcms/vcms-common/errs"
	"github.com/tars-vcms/vcms-gateway/entity/proxy"
	"github.com/tars-vcms/vcms-gateway/repo/auth"
	"github.com/tars-vcms/vcms-gateway/repo/header"
	"github.com/tars-vcms/vcms-gateway/repo/proxymanager"
	"github.com/tars-vcms/vcms-gateway/repo/routemanager"
	"net/http"
	"strconv"
)

type ProxyHttpImp struct {
	mux    *tars.TarsHttpMux
	route  routemanager.HttpRouteManager
	auth   auth.HttpAuth
	proxy  proxymanager.GatewayProxyManager
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

func (p ProxyHttpImp) handleRequest(w http.ResponseWriter, r *http.Request) {

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
	httpCode, err := gProxy.DoRequest(routeHttp, w, r)
	if err != nil {
		p.handleError(w, httpCode, err)
		return
	}
}

func newProxyHttpImp() *ProxyHttpImp {
	mux := &tars.TarsHttpMux{}
	proxyHttp := &ProxyHttpImp{
		mux:    mux,
		route:  routemanager.NewHttpRouteManager(),
		auth:   auth.NewHttpAuth(),
		proxy:  proxymanager.NewGatewayProxyManager(),
		header: header.NewHttpHeader(),
	}
	mux.HandleFunc("/", proxyHttp.handleRequest)
	return proxyHttp
}
